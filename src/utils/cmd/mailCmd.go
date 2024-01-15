package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/rbgayoivoye09/keep-online/src/utils/config"
	. "github.com/rbgayoivoye09/keep-online/src/utils/log"

	"github.com/spf13/cobra"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
)

var mailCmd = &cobra.Command{
	Use:   "mail",
	Short: "Configure keep-online settings",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.GetConfig()
		Usage(c.Mail.Name, c.Mail.Password)
	},
}

// CustomerImapClient 调用NewImapClient
func CustomerImapClient(name, password string) (*client.Client, error) {
	// 【修改】账号和密码
	return NewImapClient(name, password)
}

// NewImapClient 创建IMAP客户端
func NewImapClient(username, password string) (*client.Client, error) {
	// 【字符集】  处理us-ascii和utf-8以外的字符集(例如gbk,gb2313等)时,
	//  需要加上这行代码。
	// 【参考】 https://github.com/emersion/go-imap/wiki/Charset-handling
	imap.CharsetReader = charset.Reader

	Logger.Sugar().Info("Connecting to server...")

	// 连接邮件服务器
	c, err := client.DialTLS("imap.exmail.qq.com:993", nil)
	if err != nil {
		Logger.Sugar().Fatal(err)
	}
	Logger.Sugar().Info("Connected")

	// 使用账号密码登录
	if err := c.Login(username, password); err != nil {
		return nil, err
	}

	Logger.Sugar().Info("Logged in")

	return c, nil
}

// Usage
// 【处理业务需求】假设需求是找出求以subject开头的标题的最新邮件，并下载附件。
// 【思路】有些邮件包含附件后会变得特别大，如果要遍历的邮件很多，直接遍历处理，每封邮件都获取'RFC822'内容，
// fetch方法执行耗时可能会很长, 因此可以分两次fetch处理，减少处理时长：
// 1)第一次fetch先使用ENVELOP或者RFC822.HEADER获取邮件头信息找到满足业务需求邮件的id
// 2)第二次fetch根据这个邮件id使用'RFC822'获取邮件MIME内容，下载附件
func Usage(name, password string) {
	// 连接邮件服务器
	c, err := CustomerImapClient(name, password)
	if err != nil {
		Logger.Sugar().Fatal(err)
	}
	// Don't forget to logout
	defer c.Logout()

	// 选择收件箱
	_, err = c.Select("INBOX", false)
	if err != nil {
		Logger.Sugar().Fatal(err)
	}

	// 搜索条件实例对象
	criteria := imap.NewSearchCriteria()

	// ALL是默认条件
	// See RFC 3501 section 6.4.4 for a list of searching criteria.
	criteria.WithoutFlags = []string{"ALL"}
	ids, _ := c.Search(criteria)
	var s imap.BodySectionName

	for {
		if len(ids) == 0 {
			break
		}
		id := pop(&ids)

		seqset := new(imap.SeqSet)
		seqset.AddNum(id)
		chanMessage := make(chan *imap.Message, 1)
		go func() {
			// 第一次fetch, 只抓取邮件头，邮件标志，邮件大小等信息，执行速度快
			if err = c.Fetch(seqset,
				[]imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchRFC822Size},
				chanMessage); err != nil {
				// 【实践经验】这里遇到过的err信息是：ENVELOPE doesn't contain 10 fields
				// 原因是对方发送的邮件格式不规范，解析失败
				// 相关的issue: https://github.com/emersion/go-imap/issues/143
				Logger.Sugar().Info(seqset, err)
			}
		}()

		message := <-chanMessage
		if message == nil {
			Logger.Sugar().Info("Server didn't returned message")
			continue
		}
		Logger.Sugar().Infof("%v: %v bytes, flags=%v \n", message.SeqNum, message.Size, message.Flags)

		if strings.HasPrefix(message.Envelope.Subject, "EB VPN Password") {
			chanMsg := make(chan *imap.Message, 1)
			go func() {
				// 这里是第二次fetch, 获取邮件MIME内容
				if err = c.Fetch(seqset,
					[]imap.FetchItem{imap.FetchRFC822},
					chanMsg); err != nil {
					Logger.Sugar().Info(seqset, err)
				}
			}()

			msg := <-chanMsg
			if msg == nil {
				Logger.Sugar().Info("Server didn't returned message")
			}

			section := &s
			r := msg.GetBody(section)
			if r == nil {
				Logger.Sugar().Fatal("Server didn't returned message body")
			}

			// Create a new mail reader
			// 创建邮件阅读器
			mr, err := mail.CreateReader(r)
			if err != nil {
				Logger.Sugar().Fatal(err)
			}

			// Process each message's part
			// 处理消息体的每个part
			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					break
				} else if err != nil {
					Logger.Sugar().Fatal(err)
				}

				switch h := p.Header.(type) {
				case *mail.InlineHeader:
					// This is the message's text (can be plain-text or HTML)
					// 获取正文内容, text或者html
					b, _ := ioutil.ReadAll(p.Body)
					Logger.Sugar().Info("Got text: ", string(b))

					// 定义正则表达式模式
					pattern := `Your password: (\w+)`

					// 编译正则表达式
					re := regexp.MustCompile(pattern)

					// 查找匹配项
					matches := re.FindStringSubmatch(string(b))

					// 如果找到匹配项，则输出密码后面的字符串并写入文件
					if len(matches) > 1 {
						password := matches[1]
						Logger.Sugar().Info("Password:", password)

						// 将密码写入文件
						err := ioutil.WriteFile("password.txt", []byte(password), 0644)
						if err != nil {
							Logger.Sugar().Error("Error writing to file:", err)
						} else {
							Logger.Sugar().Info("Password written to file: password.txt")
						}
					} else {
						Logger.Sugar().Error("Password not found.")
					}
				case *mail.AttachmentHeader:
					// This is an attachment
					// 下载附件
					filename, err := h.Filename()
					if err != nil {
						Logger.Sugar().Fatal(err)
					}
					if filename != "" {
						Logger.Sugar().Info("Got attachment: ", filename)
						b, _ := ioutil.ReadAll(p.Body)
						file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
						defer file.Close()
						n, err := file.Write(b)
						if err != nil {
							Logger.Sugar().Info("写入文件异常", err.Error())
						} else {
							Logger.Sugar().Infof("写入Ok：", n)
						}
					}
				}
				Logger.Sugar().Infof("已找到满足需求的邮件")
				return
			}
		}
	}
}

func pop(list *[]uint32) uint32 {
	length := len(*list)
	lastEle := (*list)[length-1]
	*list = (*list)[:length-1]
	return lastEle
}

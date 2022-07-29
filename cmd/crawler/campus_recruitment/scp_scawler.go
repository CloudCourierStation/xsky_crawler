package campus_recruitment

import (
	"io"
	"regexp"
	"strings"
	"xsky_crawler/pkg/campus_recruitment"
	util "xsky_crawler/utils"

	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var (
	whitespaceOnly    = regexp.MustCompile("(?m)^[ \t]+$")
	leadingWhitespace = regexp.MustCompile("(?m)(^[ \t]*)(?:[^ \t\n])")
)

func NewCrawlerCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	var rootfsPath string
	cmds := &cobra.Command{
		Use:   "crawler",
		Short: "crawler",
		Long: Dedent(`
       ┌──────────────────────────────────────────────────────────┐
       │ This is xsky  crawler tools                              │
       │                                                          │
       └──────────────────────────────────────────────────────────┘
  `),
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if rootfsPath != "" {
				if err := util.Chroot(rootfsPath); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmds.AddCommand(campusRecruitmentCrawler(out))
	return cmds
}

func Dedent(text string) string {
	var margin string
	text = whitespaceOnly.ReplaceAllString(text, "")
	indents := leadingWhitespace.FindAllStringSubmatch(text, -1)
	for i, indent := range indents {
		if i == 0 {
			margin = indent[1]
		} else if strings.HasPrefix(indent[1], margin) {
			continue
		} else if strings.HasPrefix(margin, indent[1]) {
			margin = indent[1]
		} else {
			margin = ""
			break
		}
	}
	if margin != "" {
		text = regexp.MustCompile("(?m)^"+margin).ReplaceAllString(text, "")
	}
	return text
}

func campusRecruitmentCrawler(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cr",
		Short: "scratch xsky Campus Recruitment information stored in JSON files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCRCrawler(out, cmd)
		},
		Args: cobra.NoArgs,
	}
	cmd.Flags().StringP("url", "u", "https://xskydata.jobs.feishu.cn/school", "xsky campus recruitment url")
	cmd.Flags().Int("current", 1, "xsky campus recruitment page index")
	cmd.Flags().Int("limit", 100, "xsky campus recruitment page size")
	return cmd
}

func runCRCrawler(out io.Writer, cmd *cobra.Command) error {
	url, err := cmd.Flags().GetString("url")
	if err != nil {
		klog.Info(err)
		return err
	}
	current, err := cmd.Flags().GetInt("current")
	if err != nil {
		klog.Info(err)
		return err
	}
	limit, err := cmd.Flags().GetInt("limit")
	if err != nil {
		klog.Errorln(err)
		return err
	}
	err = campus_recruitment.Crawler(url,current,limit)
	if err != nil {
		klog.Errorln(err)
		return err
	}
	return nil
}

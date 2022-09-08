package markdown

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jomei/notionapi"
)

var (
	pageBuilder strings.Builder
)

func Build(page *notionapi.Page, blocks []notionapi.Block) (*string, error) {
	for _, block := range blocks {
		// bt := block.GetType()
		// fmt.Println(bt)
		switch block.GetType() {
		case notionapi.BlockTypeParagraph:
			str, err := paragraph(block)
			if err != nil {
				return nil, err
			}

			pageBuilder.WriteString(str)
		case notionapi.BlockTypeHeading1, notionapi.BlockTypeHeading2, notionapi.BlockTypeHeading3:
			str, err := heading(block)
			if err != nil {
				return nil, err
			}

			pageBuilder.WriteString(str)
		case notionapi.BlockTypeCode:
			str, err := code(page, block)
			if err != nil {
				return nil, err
			}

			pageBuilder.WriteString(*str)
		case notionapi.BlockTypeBulletedListItem:
			str, err := bulletedList(block)
			if err != nil {
				return nil, err
			}

			pageBuilder.WriteString(str)
		default:
			return nil, fmt.Errorf("currently do not support %s", block.GetType())
		}

		pageBuilder.WriteString("\n")
	}

	finalString := pageBuilder.String()

	return &finalString, nil
}

func TitleFromPage(page *notionapi.Page) string {
	pageNameProperty := page.Properties["Name"].(*notionapi.TitleProperty)
	title := pageNameProperty.Title[0].PlainText
	return title
}

func addTextModifications(richtext notionapi.RichText) string {
	// guard against non-text type of paragraphs
	text := richtext.Text.Content

	if richtext.Annotations.Bold {
		text = boldify(text)
	}

	if richtext.Annotations.Italic {
		text = italicize(text)
	}

	if richtext.Annotations.Strikethrough {
		text = strike(text)
	}

	if richtext.Annotations.Code {
		text = codify(text)
	}

	if richtext.Text.Link != nil && richtext.Text.Link.Url != "" {
		text = linkify(text, richtext.Text.Link.Url)
	}

	return text
}

func bulletedList(block notionapi.Block) (string, error) {
	blb := block.(*notionapi.BulletedListItemBlock)
	var sb strings.Builder
	sb.WriteString("* ")

	// TODO: iterate over `BulletedListItem.Children` for nested bullets
	for _, bulletBlock := range blb.BulletedListItem.RichText {
		text := addTextModifications(bulletBlock)
		sb.WriteString(text)
	}

	return sb.String(), nil
}

func paragraph(block notionapi.Block) (string, error) {
	pb := block.(*notionapi.ParagraphBlock)
	var sb strings.Builder

	if len(pb.Paragraph.RichText) == 0 {
		return "", nil
	}

	for _, textBlock := range pb.Paragraph.RichText {
		text := addTextModifications(textBlock)
		sb.WriteString(text)
	}

	return sb.String(), nil
}

func heading(block notionapi.Block) (string, error) {
	var sb strings.Builder

	sb.WriteString("\n")

	switch block.GetType() {
	case notionapi.BlockTypeHeading1:
		sb.WriteString("# ")
		b := block.(*notionapi.Heading1Block)
		sb.WriteString(b.Heading1.RichText[0].PlainText)
	case notionapi.BlockTypeHeading2:
		sb.WriteString("## ")
		b := block.(*notionapi.Heading2Block)
		sb.WriteString(b.Heading2.RichText[0].PlainText)
	case notionapi.BlockTypeHeading3:
		sb.WriteString("### ")
		b := block.(*notionapi.Heading3Block)
		sb.WriteString(b.Heading3.RichText[0].PlainText)
	}

	return sb.String(), nil
}

func code(page *notionapi.Page, block notionapi.Block) (*string, error) {
	codeBlock := block.(*notionapi.CodeBlock)
	filename := createGistFilename(page, block)
	gist, err := createOrUpdateGist(filename, codeBlock)
	if err != nil {
		return nil, err
	}

	username, err := authenticatedGithubUsername()
	if err != nil {
		return nil, err
	}

	gistString := fmt.Sprintf("{{< gist %s %s >}}", *username, *gist.ID)

	return &gistString, nil
}

func createGistFilename(page *notionapi.Page, block notionapi.Block) string {
	codeBlock := block.(*notionapi.CodeBlock)
	title := TitleFromPage(page)

	extensionsMap := map[string]string{
		"javascript": "js",
		"bash":       "sh",
		"golang":     "go",
	}

	filename := fmt.Sprintf("%s:%s.%s", ToSnakeCase(title), strings.ReplaceAll(block.GetID().String(), "-", ""), extensionsMap[codeBlock.Code.Language])

	return filename
}

func boldify(text string) string {
	return fmt.Sprintf("**%s**", text)
}

func italicize(text string) string {
	return fmt.Sprintf("*%s*", text)
}

func strike(text string) string {
	return fmt.Sprintf("~~%s~~", text)
}

func codify(text string) string {
	return fmt.Sprintf("`%s`", text)
}

func linkify(text, URL string) string {
	return fmt.Sprintf("[%s](%s)", text, URL)
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func clearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func ToSnakeCase(str string) string {
	onlyalphanumeric := clearString(str)
	spaces := strings.ReplaceAll(onlyalphanumeric, " ", "_")
	dashes := strings.ReplaceAll(spaces, "-", "_")
	return strings.ToLower(dashes)
}

# n2h — Notion to Hugo

I love writing in Notion! But want to keep my Hugo blog alive. This is my solution.

For a more comprehensive breakdown of this project, check out my blog post [Write in Notion, Publish with Hugo](https://livingissodear.com/posts/write_in_notion_publish_with_hugo_introducing_n2h/).

This package lives in your `/usr/bin` and can be called from anywhere to create a new Hugo draft from a page in Notion. Right now this supports converting this items from Notion:
* Bold
* Italic
* Strikethrough
* Inline Code
* Bulleted List (1 level deep)
* Headers
* Links
* Code Blocks (via Github Gists)

## Setup

1. Clone the repo & `cd n2h`
2. Copy secrets with `cp secrets-example.yml secrets.yml`
3. Replace variables with your data (See blog post for more info how)
4. Run `make install`

## Usage

Run `n2h --help` to see the accepted params.

If I had a blog article "How I Use Notion at Work" that lived in the database "My Thoughts" and had a "Writing Status" property that was equal to "Ready to publish", I would run:

```bash
n2h -status-property="Writing Status" -status="Ready to publish" -title="How I Use"
```

And an new `.md` file would show up in my blog directory.

Happy writing!
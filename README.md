# Dumbo

![Dumbo](https://thistestsite.space/wp-content/uploads/2018/07/071518DUMBO.jpg)

Dumbo is the world's dumbest static site generator. It does one thing - convert HTML and Markdown to HTML based on templates.

## Installation

Download Dumbo from the GitHub releases page and place it in a folder in your `PATH`, e.g. `/usr/local/bin`.

Dumbo has only been tested on Linux.

## Usage

Run Dumbo like:

```
dumbo <INPUT_DIR> <OUTPUT_DIR>
```

This will build from sources in the input directory and leave a static HTML site in the output directory. **If the output directory exists, it will be overwritten!**

### Input

Your input directory should be structured like this:

```
_templates/
  page.html
  post.html
  ...
a.html
b.post.md
z.png
c/
  d.page.md
  e/
    f.post.md
    g.html
```

Any non-Markdown files will be copied verbatim to the output directory.

Markdown files are named in the form `<name>.<template>.md`. The template must be present in the `_templates` folder. They will be rendered as HTML using this template.

Folder structure is preserved when building. For example, in the directory above, `c/e/f.post.md` will be built as `c/e/f.html`.

### Templates

Templates should be HTML files with a place to substitute a title and the Markdown content, like:

```
<html>
    <head>
        <title>
            {{.Title}}
        </title>
    </head>

    <body>
        {{.Content}}
    </body>
</html>
```

## Development

If you want to contribute, feel free!

You can build with:

```
go build
```

Then run the binary `dumbo` in the working directory to test your changes.

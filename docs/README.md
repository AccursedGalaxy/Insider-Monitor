# Solana Insider Monitor Documentation

This directory contains the documentation for the Solana Insider Monitor project, built using [Jekyll](https://jekyllrb.com/) and the [Just the Docs](https://just-the-docs.github.io/just-the-docs/) theme.

## Local Development

To build and preview the site locally:

1. Install Ruby and Bundler
2. Run `bundle install` in this directory
3. Run `bundle exec jekyll serve` to start the local server
4. Visit `http://localhost:4000/Insider-Monitor/` in your browser

## Structure

- `index.md`: The home page
- `_config.yml`: Site configuration
- `*.md`: Content pages
- `assets/`: Images, CSS, and other static files

## Adding Content

To add a new page:

1. Create a new Markdown file with front matter
2. Add appropriate navigation settings in the front matter
3. Add your content in Markdown

Example front matter:

```yaml
---
layout: default
title: My New Page
nav_order: 5
---
```

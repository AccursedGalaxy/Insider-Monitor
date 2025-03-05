# Solana Insider Monitor Documentation

This directory contains the documentation for Solana Insider Monitor, hosted on GitHub Pages. The documentation is built using Jekyll with the Just the Docs theme.

## Structure

- `_config.yml` - Jekyll configuration
- `index.md` - Home page
- `*.md` - Documentation pages
- `assets/` - Images, CSS, and other assets
- `_layouts/` - Custom layouts
- `.nojekyll` - File to prevent GitHub Pages from ignoring files that start with underscores
- `CNAME` - For custom domain setup

## Local Development

To develop or preview the documentation locally:

1. Install Ruby and Bundler
2. Install dependencies:
   ```
   bundle install
   ```
3. Start the local server:
   ```
   bundle exec jekyll serve
   ```
4. Open `http://localhost:4000` in your browser

## Adding Pages

To add a new documentation page:

1. Create a new markdown file (e.g., `new-page.md`)
2. Add front matter at the top:
   ```yaml
   ---
   layout: default
   title: Page Title
   nav_order: X
   description: Brief description of the page
   ---
   ```
3. Add the content using Markdown

## Custom Styling

Custom styles are defined in `assets/css/custom.scss`. This file customizes the Just the Docs theme with Solana Insider Monitor branding.

## Deploying

The documentation is automatically deployed to GitHub Pages when changes are pushed to the `main` branch.

## License

The documentation content is licensed under the same license as the main project.

# notion-page-visibility-checker

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
目次

- [how to set up development environment](#how-to-set-up-development-environment)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# how to set up development environment
1. clone this repository
2. run `docker compose up -d --build`
3. run `docker compose exec notion-page-visibility-checker bash`

# how to use this tool
1. fork this repository
2. set environment variables in github environment secrets
   - `NOTION_TOKEN` : your notion token
   - `DISCORD_WEBHOOK_URL` : your discord webhook url
   
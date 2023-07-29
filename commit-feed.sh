#!/bin/bash

if [ "$(git status -s)" = "" ]; then
  echo "No changes detected. Exiting."
  exit 0
fi

git config user.name "GitHub Actions Bot"
git config user.email "<>"
git add inforadio.json
git commit -m "Update feed from scheduled workflow"
git push origin main

exit 0

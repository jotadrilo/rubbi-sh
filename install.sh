#!/bin/bash

sudo cp rubbi-sh /usr/local/bin/
sudo chmod +x /usr/local/bin/rubbi-sh

echo "    Please import the .rubbi.sh file in your profile configuration or add its content directly"
echo ""
echo "     . $PWD/dotfiles/.rubbi.sh"
echo ""
echo "    or the .rubbi.minimal.sh for a minimal setup (rbsh shell function only)"
echo ""
echo "     . $PWD/dotfiles/.rubbi.minimal.sh"
echo ""
echo "    This tool is based on some alias and shell functions that are core to work."

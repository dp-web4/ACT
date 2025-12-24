#!/bin/bash
echo "Starting Society 4 Private Blockchain..."
echo "Hardware Hash: $(cat $HOME/.society4chain/hardware_binding.json | grep hardware_hash | cut -d'"' -f4)"
./society4chaind start --home $HOME/.society4chain

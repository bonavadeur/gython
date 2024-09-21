#!/bin/sh

# first, start outlier in background
python3 /ko-app/outlier/main.py &

# second, start go code
exec /ko-app/gython

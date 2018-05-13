#!/bin/sh
COMMAND="samos --data-dir $DATA_DIR --wallet-dir $WALLET_DIR $@"

adduser -D -u 10000 samos

if [[ \! -d $DATA_DIR ]]; then
    mkdir -p $DATA_DIR
fi
if [[ \! -d $WALLET_DIR ]]; then
    mkdir -p $WALLET_DIR
fi

chown -R samos:samos $( realpath $DATA_DIR )
chown -R samos:samos $( realpath $WALLET_DIR )

su samos -c "$COMMAND"

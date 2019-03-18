#!/usr/bin/env python
# -*- coding: utf-8 -*-
# File Name: fix_zerohour_bug.py
# Author: ben
# mail: zhiyuan_06@foxmail.com
# Created Time: 一  3/18 15:40:39 2019

import sys
reload(sys)
sys.setdefaultencoding("utf-8")
import commands
import json
import time

def shell_result(cmd):
   cnt=commands.getoutput(cmd)
   return json.loads(cnt)

#主账户
main_wallet = "/Users/lxx/.samos/wallets/samos_cli.wlt"
main_addr = "EX8omhDyjKtc8zHGp1KZwn7usCndaoJxSe"
#中转账户
transfer_wallet = "/Users/lxx/.samos/wallets/w1.wlt"
transfer_addr = "273K8WaBCB6kYC6W34qnfNAuoMoH1z5zpZq"

coin_process = "samos-cli"
cmd=coin_process + " addressOutputs " + main_addr
utxos=shell_result(cmd)
need_merge = False
for ux in utxos["outputs"]["head_outputs"]:
    if ux["hours"]  <= 0:
        print "has zero hour, need handle"
        print ux
        need_merge = True
        break

if need_merge:
    cmd = coin_process + " walletBalance " + main_wallet
    print cmd
    tcoins = shell_result(cmd)
    total_coins = tcoins["confirmed"]["coins"]
    ## 所有币全部发送到指定地址，等待10秒，再从指定地址转回来
    cmd = coin_process + " send -f " +  main_wallet +  " " + transfer_addr + " " + total_coins
    print "[send transfer cmd] ",  cmd
    txid = commands.getoutput(cmd)
    txid = txid.split(":")[1]
    print "[transaction id]", txid
    cmd = coin_process + " transaction " + txid 
    result = shell_result(cmd)
    print result
    while True:
        if result["transaction"]["status"]["confirmed"]:
            print "txid is confirmed"
            break
        else:
            time.sleep(2)
            result = shell_result(cmd)
    cmd = coin_process + " send -f " + transfer_wallet + " " + main_addr + " " + total_coins
    print "[send back cmd] ", cmd
    txid = commands.getoutput(cmd)
    txid = txid.split(":")[1]
    print "[transaction id]", txid

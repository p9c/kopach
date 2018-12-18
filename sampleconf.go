package main

var sampleKopachConf = `;;; URL in username:password@protocol://address:port format (default: [user:pa55word@http://127.0.0.1:11048])
;url=											
;;; RPC server certificate chain for validation, if needed (default: /home/loki/.pod/rpc.cert)
;rpccert=						
;;; Enable TLS
;tls=1       									
;;; Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)
;proxy=
;;; Username for proxy server
;proxyuser=user
;;; Password for proxy server
;proxypass=pa55word
;;; Connect to testnet
;testnet=1
;;; Connect to the simulation test network
;simnet=1
;;; Do not verify tls certificates (not recommended!)
;skipverify=1
`

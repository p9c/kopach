package main

var sampleKopachConf = `
;url=         	;;; URL in username:password@protocol://address:port format. If other ports starting at this range answer, another 8 will also be dialed (default: user:pa55word@http://127.0.0.1:11048)
;oneonly=1		;;; Only connect to the specified port
;port=			;;; specify other ports than the one in URL to connect to
;rpccert=		;;; RPC server certificate chain for validation, if needed (default: /home/loki/.pod/rpc.cert)
;tls=1			;;; Enable TLS
;proxy=			;;; Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)
;proxyuser=		;;; Username for proxy server
;proxypass=		;;; Password for proxy server
;testnet=1		;;; Connect to testnet
;simnet=1		;;; Connect to the simulation test network
;skipverify=1	;;; Do not verify tls certificates (not recommended!)
;addr=			;;; Addresses to mine to (add multiple for random selection)
;interactive=1	;;; Schedules kernel to do shorter jobs when mining with a card also running a display
`

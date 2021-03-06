module github.com/sbezverk/gobmp/pkg/message 

go 1.13

replace (
	github.com/sbezverk/gobmp/pkg/base => ../base
	github.com/sbezverk/gobmp/pkg/bgp => ../bgp
	github.com/sbezverk/gobmp/pkg/bgpls => ../bgpls
	github.com/sbezverk/gobmp/pkg/bmp => ../bmp
	github.com/sbezverk/gobmp/pkg/gobmpsrv => ../gobmpsrv
	github.com/sbezverk/gobmp/pkg/kafka => ../kafka
	github.com/sbezverk/gobmp/pkg/ls => ../ls
	github.com/sbezverk/gobmp/pkg/message => ../message
	github.com/sbezverk/gobmp/pkg/parser => ../parser
	github.com/sbezverk/gobmp/pkg/pub => ../pub
	github.com/sbezverk/gobmp/pkg/sr => ../sr
	github.com/sbezverk/gobmp/pkg/srv6 => ../srv6
	github.com/sbezverk/gobmp/pkg/tools => ../tools
)

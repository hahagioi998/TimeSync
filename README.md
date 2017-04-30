#TimeSync

use:

	server: TimeSync -type server -ServerAddr 192.168.123.10:2345
	client:	TimeSync -type client -ServerAddr 255.255.255.255:2345 -Client 192.168.123.11:0

# pex-cmd
Command line tools for automated PEX tasks 
`pex-cmd [command]`
o
`go run main.go [command]`
o
`docker-compose run cmd go run main.go [command]`

## Command mf
mf Operation tool for Merchant files

### Check merchant files on the FTP server
`mf download list`


### Download an individual FTP server file
`mf download file [filename]`


### Download all .gz files from the FTP server
`mf download all`


### Check merchants files .xml available to upload to the database
`mf upload list`


### Upload merchant file to the database
`mf upload file [filename]`


### Upload all merchant file .xml files to the database
`mf upload all`

### Verbose Mode
The verbose mode allows to visualize the errors and successes when executing a command, At the moment only available for `mf upload file` and `mf upload all`.

`mf upload file -v [filename]`
`mf upload all -v`

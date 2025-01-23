## Instructions on bulding the code
#### It is tested with golang version 1.22.5
```
cd command_line_tool
go mod init command_line_tool
go build
```
## Instructions on running the application
```
# output the command line options
./command_line_tool -help

# validate the input task list and output the expected total runtime without running the tasks
./command_line_tool -taskfile task_list.txt -dryrun

# run the tasks and determine the difference in the actual runtime versus the expected runtime
./command_line_tool -taskfile task_list.txt -difftime

# validate the imput task file containing errors
./command_line_tool -taskfile task_list_error.txt -dryrun
```

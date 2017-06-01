# cobra

## Purpose

When I write android frontend for our project,  our
backend always be unavailable because of the modification to db
and incomplete Regression Testing.

So, I write this tool to monitor whether service is unavailable,
and notify the backend developer if it is.

## How to use
There are two types of file in this tool. The first one is config
file which is used to config the tool, the second one is rule file
which contains url, method, params of interface you wanna test.

Now we just support HTTP.

See **example** folder for detail of the format of these two files.

Run the following command to start tool when you have finished the config file and rule:

**cobra -f config_file_path**

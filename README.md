下载地址  
`http://www.oracle.com/technetwork/database/database-technologies/instant-client/downloads/index.html`

选择：  
`Instant Client for Linux x86-64`  

1）下载以下ZIP包：

`instantclient-basic-linux.x64-12.2.0.1.0.zip 已经在项目目录添加，可以在适应也可以自己下载`

2）安装依赖包：  
`yum install libaio* -y`  
3）上传到服务器并解压到相同目录下  

``` 
unzip instantclient-basic-linux.x64-12.2.0.1.0.zip  
unzip instantclient-sqlplus-linux.x64-12.2.0.1.0.zip  
sudo mv instantclient_12_2 /soft/

```   
4）添加环境变量  

``` 
export ORACLE_HOME=/home/go/instantclient_12_2  
export LD_LIBRARY_PATH=/home/go/instantclient_12_2  
export TNS_ADMIN=/home/go/instantclient_12_2  
export PATH=$PATH:/home/go/instantclient_12_2  `

``` 

      
5）报错处理     

``` 
./oracle2xlsx 
panic: ORA-00000: DPI-1047: Cannot locate a 64-bit Oracle Client library: "libclntsh.so: cannot open shared object file: No such file or directory". See https://oracle.github.io/odpi/doc/installation.html#linux for hel`

解决：
`vi  /etc/ld.so.conf  
include ld.so.conf.d/*.conf
/home/go/instantclient_12_2`

``` 
  
  
运行  
`./oracle2xlsx -h 127.0.0.1 -u oracle -p oracle -s test1 -t ./user.xlsx`

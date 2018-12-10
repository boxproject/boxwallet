### 概念
此处是Import package方式引用boxwallet的入口，startup包含了启动依赖项，修改这个go file，可以自定义配置文件目录等，然后这个包不能和mock层同时import，2个是并列层级的，相当于一个是沙盒，一个是生产环境的概念，2个不能同时存在
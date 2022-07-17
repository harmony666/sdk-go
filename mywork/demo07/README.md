# demo07文件说明
- 其中`tmp`文件夹是使用fresh实现热加载生成的文件夹 https://github.com/gravityblast/fresh
- 简单实现对从客户端发来的post请求中携带的json进行schema校验,并进行简单权限控制。
- 通过使用json path来获取客户端传来的json中`datatype`这个键所对应的值，然后通过这个值与配置文件`config.yaml`中的`name`字段进行比较，然后找到schema文件路径。
- post请求头中携带有base64编码的证书，在服务端通过x509对这个证书进行解析，拿到其中的`Organization`字段，然后通过读取的组织信息和配置文件`config.yaml`中的字段`Role`进行比较，查找配置文件中是否有该组织，如果有该组织那么检查该组织是否有权限访问该接口。
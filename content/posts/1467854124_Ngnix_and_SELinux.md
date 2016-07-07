
::title:Nginx and SELinux
::timestamp:1467854124

# Nginx and SELinux

I am learning nginx and decided to use a RHEL linux server to do so. I installed and ran nginx and hit the delightful "Welcome to nginx!" when accessing the public IP (e.g. 1.2.3.4) of the site. Okay, all good.

I wanted to use server_name to run multiple webservers on the same port on the same machine. I pointed a different URL of mine at the IP (1.2.3.4) using (e.g. example.com) and created a second config file that would address example.com and serve a very basic HTML file from a different directory.

```
server {
    listen 80;
    server_name example.com;
    root /var/www/example.com;

    location / {
        try_files $uri /index.html;
    }
}
```

I then modified the default.conf to only handle the IP of the site.

```
server {
    listen       80;
    server_name  1.2.3.4;

    location / {
        root   /usr/share/nginx/html;
        index  index.html index.htm;
    }

    error_page   500 502 503 504  /50x.html;
    location = /50x.html {
        root   /usr/share/nginx/html;
    }
}
```

Okay, looks good. I restarted nginx and deciding to check 1.2.3.4. It showed me the welcome page as expected. However, example.com showed a 403 forbidden error.

Let's check the permissions:

```bash
$ namei -l /var/www/example.com/index.html
f: /var/www/example.com/index.html
dr-xr-xr-x root root /
drwxr-xr-x root root var
drwxr-xr-x root root www
drwxr-xr-x root root example.com
-rw-r--r-- root root index.html
```

They look fine to me... so what's going on? Let's check to see if it matches the directory nginx is currently serving out of for 1.2.3.4.

```bash
/usr/share/nginx/html/index.html
dr-xr-xr-x root root /
drwxr-xr-x root root usr
drwxr-xr-x root root share
drwxr-xr-x root root nginx
drwxr-xr-x root root html
-rw-r--r-- root root index.html
```

They seem to match to me. What is going on... I found the error log for nginx under /var/log/nginx/error.log which showed me this line:

```
... "/var/www/example.com/index.html" is forbidden (13: Permission denied)...
```

Not very helpful, after some research I found [this post](https://stackoverflow.com/questions/22586166/why-does-nginx-return-a-403-even-though-all-permissions-are-set-properly) which fixed my problem. It turns out that this thing called SELinux was stopping nginx from accessing that directory. I had never heard of SELinux before and after some research I found that it was a very powerful security tool for managing fine-grained access between processes and objects like ports, files, and dirs.  It does this by tagging objects with a user, role, and type. The type is most important here. A process tagged with type A can only access certain types. For example, we can see nginx's type using the Z flag.

```
$ ps -eZ | grep nginx
system_u:system_r:httpd_t:s0    30956 ?        00:00:00 nginx
system_u:system_r:httpd_t:s0    30957 ?        00:00:00 nginx
```

The Z flag allows us to see the security data in the from of user:role:type. You can see that the type of the nginx process is httpd_t. It can access certain types like usr_t and httpd_sys_content_t. The content I was attempting to serve was var_t, something httpd_t processes can not access. All I had to do was change the type of my files that I am serving to a proper type and everything should be dandy. I will also changed the user to be system_u to make things consistent.

```bash
$ chcon -R -t httpd_sys_content_t -u system_u /var/www
$ ls -Zd /var/www
drwxr-xr-x. root root system_u:object_r:httpd_sys_content_t:s0 www/
$ ls -Zd /var/www/example.com
drwxr-xr-x. root root system_u:object_r:httpd_sys_content_t:s0 example.com/
$ ls -Zd /var/www/example.com/index.html
-rw-r--r--. root root system_u:object_r:httpd_sys_content_t:s0 index.html
```

Now, I restarted nginx for good measure and I can now see my page with 200 OK.

I'd highly recommened watching [this video](https://www.youtube.com/watch?v=MxjenQ31b70) and learning some things about SELinux if your a webmaster, web-admin, using Apache, NGINX, httpd or similar. Some useful commands to checkout are:

```
getsebool -- gets a list of all booleans SELinux has
setsebool -- sets a boolean, temporarily or permanently
chcon -- change the user, role, or type of a dir/file
sestatus -- gets the current status of SELinux
getenforce -- gets the current enforcing mode of SELinux
setenforce -- change the current enforcing mode, (i.e. turns SELinux on/off)

Contains the booleans that differ from default:
/etc/selinux/targeted/modules/active

Contains the default modes:
/etc/sysconfig/selinux
```

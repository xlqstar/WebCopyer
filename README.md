#基础使用

##简介

初次使用golang，匆忙之作，代码不敢称优雅，仅为了给我的just项目批量制作博客模板而半路写的这个webcopyer，有bug还请issue或pull request。

    前两天群里有朋友提出了两个比较犀利的问题：

    1、遇到防爬虫的情况怎么办，
    2、wget一条命令就能搞定的事情何必这么折腾。

    是这样的：

    首先，我暂时没考虑防爬虫问题，但理论上讲，只要是你人通过浏览器访问到的页面，都可以跑到，但这不是我关注的点（我没这需求），也没精力去突破这个问题。
    其次，相比较于wget或者其他类似工具而言，有这样几个有点，1、完全开源，代码思路清晰，方便二次开发。2、能自定义站点结构，免去了不必要的麻烦。3、可以下载html关联的css中再关联的资源，这个功能确保了webcopyer的真正实用性。


##使用方法及其应用场景

注意：进行以下操作时请确保配置文件存在且配置正常，配置相关信息请前往：

**一、拷贝`[ http://lyric.im ]`的模板，包括其所有资源（css、js、images）以及资源关联的资源,操作命令如下：**

    webcopyer http://lyric.im/
    或者（两条命令皆可）：
    webcopyer get http://lyric.im/
PS:这条命令是主要命令，可以说百分之百之九十的情况下用此命令即可，剩下的命令都是一些特殊应用场景下才会用到

**二、拷贝本地模板，获取其资源以及跟资源关联的资源**

    webcopyer getLocal F:\template\index.html  http://drizzlep.diandian.com/page/2/

PS:这条命令的应用场景比较特殊，因为我遇到过这样一种情况，某个博客他不分页，几百篇文章在一个页面里，导致整个网页很大，还下载了很多无用的图片等资源，我后期删选很不方便，而我仅仅拷贝其模版而已，对我来说两篇日志一页跟三篇日志一页没有任何区别，于是我就考虑这样来操作：打开该网页，查看其源码，全选、复制，将其所有html代码原样拷贝至本地文件（如：F:\template\index.html），然后我先根据需要进行编辑，再用webcopyer来分析这个修改过后的本地html，就太随意啦！~不要要注意命令后面还要带上该页面在网络上原本的url地址（如：[ http://drizzlep.diandian.com/page/2/ ]），这是内部的实现机制决定的，有这个地址才能推算出那些关联资源的正确下载路径以及一些细小的其他原因，总之用getLocal的时候记得带上这个参数！

我想这个命令的应用场景可能不止如此，待您自己发掘吧！

**三、仅拷贝html，及其直接关联资源**

    webcopyer getHtml http://lyric.im/

PS:getHtml 跟 get 的两个命令的区别是，get不仅仅获取html，以及html包含的css、js等，还包括了其包含的css中包含的其他资源，也就是下载html及其资源及其资源包含的资源，可能有点绕，语言好难表述，我想你应该懂我的意思。而getHtml仅下载html及其直接包含的css、js等资源。应该说这条命令的应用场景更少了，反正我没有要用到过，设置这条命令的原因仅仅因为反正这是我代码中的一个子函数，就索性开发出来，看各自有没有可能的应用场景需要。

**四、仅拷css，及其直接关联资源**

    webcopyer getHtml http://lyric.im/styles/style.css
    
PS:这条命令的完整描述应该是：仅仅下载css及css中包含的资源，理解第上一条命令的话这条命令应该也能秒懂了。应用场景不详……

##配置
默认情况下，请确保webcopyer同目录下有名为config文本文件，并且配置无误。倘若想设定其他路径的配置文件请使用如：[ webcopyer http://www.baidu.com --config=F:\template\config ]这样的命令，即带上[ --config ]。

我的默认配置如下：
    
    #生成目录
    destDir:    F:\kuaipan\Projects\webcopyer\template
    
    #html文件的存储路径（相对于destDir,`\`即等于destDir）
    html_dir:    \
    
    #图片资源后缀
    img_ext:	.jpg|.gif|.png|.jpeg|.ico|.JPG|.GIF|.PNG|.JPEG|.ICO
    #图片资源下载目录（相对于destDir）
    img_dir:	\theme\images
    
    #css资源后缀
    css_ext:	.css|.less|.CSS|.LESS
    #css资源下载目录（相对于destDir）
    css_dir:	\theme\css
    
    #Js资源后缀
    js_ext:		.js|.JS
    #js资源下载目录（相对于destDir）
    js_dir:		\theme\js
    
    #其他资源后缀
    other_ext:	.dll|.DLL
    #其他资源下载目录（相对于destDir）
    other_dir:	\theme\other
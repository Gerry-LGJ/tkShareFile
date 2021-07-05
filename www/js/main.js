
function createHttpRequest()
{
    var xmlHttp=null;
    try{
        // Firefox, Opera 8.0+, Safari
        xmlHttp=new XMLHttpRequest();
    }catch (e){
        // Internet Explorer
        try{
            xmlHttp=new ActiveXObject("Msxml2.XMLHTTP");
        }catch (e){
            try{
                xmlHttp=new ActiveXObject("Microsoft.XMLHTTP");
            }catch (e){
                alert("您的浏览器不支持AJAX！");
            }
        }
    }
    return xmlHttp;
}
/**
 * 本人是前端菜鸟，其实可以直接用submit，这样写只为了有个好看的进度条而已,555.
 */
function uploadFileToServer() {
    var uploadTip = document.getElementById("upload_tip_id");           //捕获提示标签
    var uploadProgress = document.getElementById("upload_progress_id"); //捕获进度条标签
    var file_obj = document.getElementById("upload_file_id").files[0];  //捕获表单中的文件标签
    uploadProgress.style.visibility = "visible";                        //显示进度条标签
    if(file_obj){                                                       //检测是否有文件被选择
        var url = "/upload";
        var form = new FormData(); 
        form.append("file", file_obj);
        var xhr = new XMLHttpRequest();
        xhr.onload = function(e) {                                      // 添加上传成功后的回调函数
            uploadTip.innerText = "上传成功";
            uploadProgress.style.visibility = "hidden";                 // 上传后隐藏进度条
            uploadProgress.value = 0;                                   // 进度条归零
        };
        xhr.onerror =  function(e){ uploadTip.innerText = "上传失败，或者您没有上传权限，请联系管理员。"; }; // 上传失败后的回调函数
        xhr.upload.onprogress = function(e) { uploadProgress.value = e.loaded*100/e.total;}; // 添加 监听上传进度函数
        xhr.open("POST", url, true);
        xhr.send(form);
    }else{
        uploadTip.innerText="请先选择文件后再上传";
    }
}

function showUploadDialog(){
    var up_dialog = document.getElementById("upload_dialog");   // 捕获上传窗口的标签
    document.getElementById("upload_tip_id").innerText = "请选择要上传的文件";
    document.getElementById("upload_progress_id").style.visibility = "hidden"; // 先隐藏上传进度条
    up_dialog.style.visibility = "visible";                     // 显示上传窗口
}
function hideUploadDialog(){
    var up_dialog = document.getElementById("upload_dialog");   // 捕获上传窗口的标签
    document.getElementById("upload_progress_id").style.visibility = "hidden"; // 隐藏上传进度条
    up_dialog.style.visibility = "hidden";                      // 隐藏上传窗口
}
/**
 * 为媒体文件添加试听或试看的功能
 * @param {后端提供按钮序号} i 方便索引媒体文件的src而已，因为不会通过button索引兄弟标签<a>里的href,555.
 */
function player(i) {
    var objs = document.getElementsByClassName("play_btn");     //捕获所有播放按钮标签
    var a = objs[i].previousElementSibling;                     //捕获每个按钮标签的前一个兄弟标签
    //新建一个窗口
    win = window.open("/player?src=" + a.href, "newW",
        "height=500,width=880,top=150,left=250,location=no,menubar=no,toolbar=no,status=no",
        true                                                    // 替换浏览器的历史纪录
    );
}
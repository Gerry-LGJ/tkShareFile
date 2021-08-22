
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
    var upload_dir = document.getElementById("upload_dir");             //捕获当前路径的span标签
    var uploadTip = document.getElementById("upload_tip_id");           //捕获提示标签
    var uploadProgress = document.getElementById("upload_progress_id"); //捕获进度条标签
    var file_obj = document.getElementById("upload_file_id").files[0];  //捕获表单中的文件标签
    uploadProgress.style.visibility = "visible";                        //显示进度条标签
    if(file_obj){                                                       //检测是否有文件被选择
        var url = "/upload";
        var form = new FormData();
        form.append("file", file_obj);                                  //添加上传文件
        var xhr = new XMLHttpRequest();
        xhr.onload = function(e) {                                      // 添加上传成功后的回调函数
            if(xhr.readyState == 4 && xhr.status == 200) {
                var responseObject = JSON.parse(xhr.responseText);
                switch(responseObject["status"])
                {
                case 0:
                    uploadTip.innerText = responseObject["msg"];
                    uploadProgress.style.visibility = "hidden";                 // 上传后隐藏进度条
                    uploadProgress.value = 0;                                   // 进度条归零
                    break;
                case 1:
                    uploadTip.innerText = responseObject["msg"];
                    uploadProgress.style.visibility = "hidden";
                    uploadProgress.value = 0;
                    break;
                default:
                    break;
                }
            } else {
                uploadTip.innerText = "上传失败，或者您没有上传权限，请联系管理员。";
            }
        };
        xhr.onerror =  function(e){ uploadTip.innerText = "上传失败，或者您没有上传权限，请联系管理员。"; }; // 上传失败后的回调函数
        xhr.upload.onprogress = function(e) { uploadProgress.value = e.loaded*100/e.total;}; // 添加 监听上传进度函数
        xhr.open("POST", url, true);
        xhr.setRequestHeader("upload-dir", upload_dir.innerText);       //把指定上传的路径放在请求头里面
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

//双击或单击a标签
var clickNum = 0;
var clickTimeId;
function aTagOnClick(url) {
    clickNum++;
    if (clickNum == 2) {
        clickNum = 0;
        clearTimeout(clickTimeId);
        //新建一个窗口
        /*window.open("/preview?src=http://"+ String(window.location.host) + url, "newW",
            "height=500,width=880,top=150,left=250,location=no,menubar=no,toolbar=no,status=no",
            //true                                                    // 替换浏览器的历史纪录
        );*/
        window.open("/preview?src=http://"+ String(window.location.host) + url);
        return
    }
    clickTimeId = setTimeout(function () {
        clickNum = 0;
        var a = document.createElement("a");
        a.style.visibility = "hidden";
        $("body").append(a);                                        // 修复firefox中无法触发click
        a.href = url;
        a.click();
        $(a).remove();
    }, 500);                                                        // 单击第一次在500ms内再单击一次就算双击
}

function aTagOnMouseOver(tip_str) {
    var onmouse_tip_id = document.getElementById("onmouse_tip_id");
    onmouse_tip_id.innerText = String(tip_str);
}
function aTagOnMouseOut() {
    var onmouse_tip_id = document.getElementById("onmouse_tip_id");
    onmouse_tip_id.innerText = "";
}

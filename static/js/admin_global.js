/* 后台全局http错误处理 */
$(document).ajaxError(function(event, jqxhr){
    console.log(event, jqxhr)
    if (jqxhr.status == 401) {
        alert(jqxhr.responseJSON.msg);
        window.location.href = "/login";
    }
});

/* 全局弹窗 */
function pop_window_sm_center(title, content) {
    $("#myModal_sm_center .modal-title").text(title);
    $("#myModal_sm_center .modal-body").text(content);
    $("#myModal_sm_center").modal("show");
}

function pop_window_sm_center_close() {
    $("#myModal_sm_center .modal-title").text('');
    $("#myModal_sm_center .modal-body").text('');
    $("#myModal_sm_center").modal("hide");
}

function pop_window_remind(title, content) {
    $("#myModal_remind .modal-title").text(title);
    $("#myModal_remind .modal-body").text(content);
    $("#myModal_remind").modal("show");
}

function pop_window_remind_close() {
    $("#myModal_remind .modal-title").text('');
    $("#myModal_remind .modal-body").text('');
    $("#myModal_remind").modal("hide");
}
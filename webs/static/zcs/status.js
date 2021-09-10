
function formatSize(size){
	if(size < 1000){
		return size + "Byte";
	}
	size /= 1024;
	if(size < 1000){
		return size.toFixed(2) + "KB";
	}
	size /= 1024;
	if(size < 1000){
		return size.toFixed(2) + "MB";
	}
	size /= 1024;
	if(size < 1000){
		return size.toFixed(2) + "GB";
	}
	size /= 1024;
	return size.toFixed(2) + "TB";
}

function formatCpuTime(time){
	time = Math.floor(time * 1000);
	const millisecond = Math.floor(time % 1000); time /= 1000;
	const second = Math.floor(time % 60); time /= 60;
	const minute = Math.floor(time % 60); time /= 60;
	const hour = Math.floor(time % 24); time /= 24;
	const day = Math.floor(time);
	var string = second + '.' + millisecond;
	if(minute){
		string = minute + ':' + string;
	}
	if(hour){
		string = hour + ':' + string;
	}
	if(day){
		string = day + ':' + string;
	}
	return string;
}

var svrid = "";
var flushtimeout = 1000;

$(document).ready(function(){
	function flushInfo(){
		$.ajax({
			url: "/zcs/api/svrinfo?svr=" + SERVER_NAME,
			type: "GET",
			success: function(res){
				svrid = res.id;
				flushtimeout = Math.max(res.interval * 1000, 1000);
				$("#status-cpu_num-v").text(res.cpu_num);
				$("#status-cpu_num").show();
				$("#status-java_version-v").text(res.java_version);
				$("#status-java_version").show();
				$("#status-os-v").text(res.os);
				$("#status-os").show();
				$("#status-max_mem-v").text(formatSize(res.max_mem));
				$("#status-max_mem").show();
			}
		});
	}
	var flushing_status = false;
	function flushStatus(){
		if(flushing_status){
			return;
		}
		flushing_status = true;
		$.ajax({
			url: "/zcs/api/svrstatus?svr=" + SERVER_NAME,
			type: "GET",
			success: function(res){
				const rid = res.id;
				if(rid != svrid){
					flushInfo();
				}
				const errstr = res.errstr;
				if(errstr){
					$("#status-errstr-v").text(res.errstr);
					$("#status-errstr").show();
					return;
				}
				$("#status-status-v").text(res.status);
				$("#status-status").show();
				$("#status-ticks-v").text(res.ticks);
				$("#status-ticks").show();
				$("#status-total_mem-v").text(formatSize(res.total_mem));
				$("#status-total_mem").show();
				$("#status-used_mem-v").text(formatSize(res.used_mem));
				$("#status-used_mem").show();
				$("#status-cpu_load-v").text(res.cpu_load.toFixed(2));
				$("#status-cpu_load").show();
				$("#status-cpu_time-v").text(formatCpuTime(res.cpu_time));
				$("#status-cpu_time").show();
			},
			complete: function(){
				flushing_status = false;
			}
		});
	}
	function timeFlushStatus(){
		flushStatus();
		setTimeout(timeFlushStatus, flushtimeout);
	}
	timeFlushStatus();
});


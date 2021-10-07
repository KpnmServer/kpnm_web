
'use strict';

const AJAX = new XMLHttpRequest();
var TARGET = "";

self.onmessage = function(event){
	const res = event.data;
	TARGET = self.origin + "/server/" + res + "/status";
	timeFlushStatus();
	self.onmessage = null;
}

AJAX.onreadystatechange = function(){
	if(AJAX.readyState === 4){
		if(AJAX.status === 200){
			postMessage(JSON.parse(AJAX.responseText));
		}else{
			postMessage({
				"status": "error",
				"error": "ServerStatusError",
				"errorMessage": AJAX.statusText
			});
		}
	}else{
		// postMessage({
		// 	"status": "error",
		// 	"error": "&LoadingError"
		// });
	}
}

AJAX.onerror = function(){
	postMessage({
		"status": "error",
		"error": "ServerStatusError",
		"errorMessage": AJAX.statusText
	});
}

function flushStatus(){
	AJAX.open('GET', TARGET, false);
	AJAX.send(null);
}
function timeFlushStatus(){
	flushStatus();
	setTimeout(timeFlushStatus, 5000);
}


"use strict";

$(document).ready(function(){
	if(window.Worker){
		let worker = new window.Worker('/static/server/infogetter.js');
		worker.onmessage = function (event){
			const res = event.data;
			if(res.status === "ok"){
				$("#status-error").hide();
				$("#status-box").show();
				$("#status-favicon").prop("src", res.favicon);
				$("#status-favicon").show();
				$("#status-desc").text(res.desc);
				$("#status-ip").text(res.ip);
				$("#status-version").text(res.version);
				$("#status-ping").text(res.ping);
				$("#status-player_count").text(res.player_count);
				$("#status-player_max_count").text(res.player_max_count);
				$("#status-players").html("");
				res.players.sort((a, b)=>a.name.localeCompare(b.name)).forEach((player)=>{
					$("#status-players").append($("<li></li>").addClass("status-players-item").text(player.name));
				});
				return;
			}
			if(res.status === "error"){
				$("#status-box").hide();
				if(res.error === "&LoadingError"){
					$("#status-error-text").hide();
					$("#status-error-loading").show();
				}else{
					$("#status-error-loading").hide();
					$("#status-error-text").text(res.errorMessage);
					$("#status-error-text").show();
				}
				$("#status-error").show();
			}
		}
		worker.postMessage(SERVER_NAME);
	}else{
		function flushStatus(){
			$.ajax({
				async: false,
				url: "/server/" + SERVER_NAME + "/status",
				type: "GET",
				success: function(res){
					if(res.status === "ok"){
						$("#status-error").hide();
						$("#status-box").show();
						$("#status-favicon").prop("src", res.favicon);
						$("#status-favicon").show();
						$("#status-desc").text(res.desc);
						$("#status-ip").text(res.ip);
						$("#status-version").text(res.version);
						$("#status-ping").text(res.ping);
						$("#status-player_count").text(res.player_count);
						$("#status-player_max_count").text(res.player_max_count);
						$("#status-players").html("");
						res.players.sort((a, b)=>a.name.localeCompare(b.name)).forEach((player)=>{
							$("#status-players").append($("<li></li>").addClass("status-players-item").text(player.name));
						});
						return;
					}
					if(res.status === "error"){
						$("#status-box").hide();
						$("#status-error-loading").hide();
						$("#status-error-text").text(res.errorMessage);
						$("#status-error-text").show();
						$("#status-error").show();
					}
				},
				failed: function(res){
					$("#status-box").hide();
					$("#status-error").text("" + res + typeof res);
					$("#status-error").show();
				}
			});
		}
		function timeFlushStatus(){
			flushStatus();
			setTimeout(timeFlushStatus, 5000);
		}
		setTimeout(timeFlushStatus, 200);
	}
});


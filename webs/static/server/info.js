
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
				$("#status-error").text(res.errorMessage);
				$("#status-error").show();
				// $("#status-favicon").hide();
				// $("#status-desc").text("None");
				// $("#status-ip").text("None");
				// $("#status-version").text("None");
				// $("#status-ping").text("None");
				// $("#status-player_count").text("None");
				// $("#status-player_max_count").text("None");
				// $("#status-players").html("");
			}
		}
	});
}

$(document).ready(function(){
	function timeFlushStatus(){
		flushStatus();
		setTimeout(timeFlushStatus, 5000);
	}
	document.getElementById("info-readme").attachShadow({mode: "open"}).appendChild(
		$(`<div id="body">`).load("/server/" + SERVER_NAME + "/infome")[0]);
	setTimeout(timeFlushStatus, 500);
});


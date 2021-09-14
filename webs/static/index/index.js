
"use strict";

$(document).ready(function(){
	var userdescDOM = $("#header-user-desc");
	userdescDOM.hide();
	$("#header-user").hover(function(){
		userdescDOM.fadeToggle(200);
	}, function(){
		userdescDOM.fadeOut(200);
	})
	$("#header-user-desc-logout-btn").click(function(){
		if(!confirm("Are you sure to log out?")){
			return;
		}
		$.ajax({
			url: "/user/api/logout",
			type: "POST",
			success: function(res){
				window.location.reload();
			}
		})
	})
})

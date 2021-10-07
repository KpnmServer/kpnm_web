
'use strict';

(async function(){
	await loadJsAsync('/static/lib/loading.js');
	function setup($){
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
	}
	loading({
		reslove: true,
		call: async function(rs){
			await loadJsAsync('/static/lib/jquery-3.6.0.min.js');
			$(document).ready(function(){
				setup($);
				rs();
			});
		}
	});
})();

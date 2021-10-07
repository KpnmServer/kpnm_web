
'use strict';

(async function(){
	await loadJsAsync('/static/lib/loading.js');

	function setup($){
		$.ajax({
			url: "/zcs/api/getcycleimgs",
			type: "GET",
			success: function(res){
				//
			}
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

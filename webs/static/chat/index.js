
'use strict';

(async function(){
	await loadJsAsync('/static/lib/loading.js');
	function setup($){
		function focuschange(id=''){
			if(id.length > 0 && id[0] === '#'){
				id = id.substring(1);
			}
			if(id === ''){
				console.log("back to white");
				return;
			}
			console.log("force id:", id);
		}
		$(window).on('hashchange', function(){
			focuschange(window.location.hash);
		});
		focuschange(window.location.hash);
	}
	loading({
		reslove: true,
		call: async function(rs){
			var st = new Date().getTime();
			await loadJsAsync('/static/lib/jquery-3.6.0.min.js');
			$(document).ready(function(){
				setup($);
				let ut = new Date().getTime() - st;
				setTimeout(()=>{rs()}, Math.max(1500 - ut, 0));
			});
		}
	});
})();

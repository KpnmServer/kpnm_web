
'use strict';

(function(){
	const body = document.getElementsByTagName('body')[0];
	function loadJs(src){
		if(document.querySelector('script[src="' + src + '"]') === null){
			var sc = document.createElement('script');
			sc.type = 'text/javascript';
			sc.src = src;
			body.appendChild(sc);
			return sc;
		}
		return null;
	}
	window.loadJs = loadJs;
	window.loadJsList = function(srclist){
		var sclist = [];
		srclist.forEach((src)=>{
			sclist.push(loadJs(src));
		});
		return sclist;
	}
	window.loadJsAsync = async function(src){
		var sc = loadJs(src);
		if(sc === null){
			return null;
		}
		return await new Promise((rs)=>{
			sc.onload = function(){
				rs(sc);
			}
		});
	}
	window.loadJsListAsync = async function(srclist){
		var sclist = [];
		srclist.forEach((src)=>{
			sclist.push(loadJsAsync(src));
		});
		return await Promise.all(sclist);
	}
})();

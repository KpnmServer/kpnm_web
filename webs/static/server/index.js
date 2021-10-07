
'use strict';

(async function(){
	await loadJsAsync('/static/lib/loading.js');
	function setup($){
		function search(keyword, force=false){
			if(!force && loading.isloading()){ return; }
			loading({
				reslove: true,
				call: function(rs){
					if(!keyword.length){
						keyword = '%'
					}else if(keyword[0] === '+'){
						keyword = keyword.substring(1);
					}else{
						keyword = '%' + keyword.replace(/[%_]/g, (s)=>{return ('\\' + s)}) + '%';
					}
					$.ajax({
						url: '/server/api/search?key=' + escape(keyword),
						type: 'GET',
						success: function(res){
							if(res.status === "ok"){
								$SERVER_PAGE.empty();
								if(!res.data.length){
									$SERVER_PAGE.html(`No content found`);
									return;
								}
								res.data.forEach((id)=>{
									$.ajax({
										async: false,
										url: '/server/api/info?id=' + id,
										type: 'GET',
										success: function(res){
											let item = $(`<div class="server-item">
	<h4 class="server-item-name"></h4>
	<i>version: <span class="server-item-version"></span></i><br/>
	<p class="server-item-desc"></p></div>`);
											item.find('.server-item-name').html($(`<a href="/server/${res.data.id}"></a>`).text(res.data.name)).
												append($(`<i style="font-size: 0.7em;"></i>`).text('(' + res.data.id + ')'));
											item.find('.server-item-version').text(res.data.version);
											item.find('.server-item-desc').text(res.data.desc);
											$SERVER_PAGE.append(item);
										}
									});
								})
								return;
							}
						},
						complete: function(){
							rs();
						}
					});
				}
			});
		}
		const $SERVER_PAGE = $('#server-page');
		{
			let keyword = /(?:key=)([^&]*)/.exec(window.location.search);
			keyword = unescape(keyword ?keyword[1] :'');
			$('#search-text').val(keyword);
			search(keyword, true);
		}
		$('#search-submit').click(function(){
			const val = $('#search-text').val();
			window.history.replaceState(null, null, window.location.pathname + (val.length ?'?key=' + escape(val) :''));
			search(val);
		});
		$('#search-text').keyup(function(event){ if(event.keyCode == 13){$('#search-submit').click();} });
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

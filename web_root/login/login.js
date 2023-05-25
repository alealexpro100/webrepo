(function(e,t) {
    window.addEventListener('load', init);
    function init() {
        $('#login').on('click', function() {
            var username=$("#username").val();
            var password=$("#password").val();
            window.location.href = "/api/login?username="+username+"&password="+password;
        });
    }
})();
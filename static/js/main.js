var oc = oc || {};

$(document).ready(function() {
    $('#search-form').submit(function(e) {
        e.preventDefault();
        var $searchText = $('#search-text');
        var searchText = $searchText.val();
        $searchText.val("");
        oc.search.page(searchText, '#search-text', true, null);
    });
    $('#login-form').submit(function(e) {
        e.preventDefault();
        var userName = $('#login-uname').val();
        var password = $('#login-pwd').val();
        oc.user.login(userName, password);
    });
    $('#reg-form').submit(function(e) {
        e.preventDefault();
        var userName = $('#reg-uname').val();
        var password = $('#reg-pwd').val();
        var passwordDuplicated = $('#reg-pwd-dup').val();
        oc.user.register(userName, password, passwordDuplicated);
    });
    $('#logout-form').submit(function(e) {
        e.preventDefault();
        oc.user.logout();
    });
    $('#home-button').on('click', function() {
        window.location.href = "/";
    });
});
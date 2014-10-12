oc = oc || {};

oc.user = oc.user || {};

oc.user.login = function(username, password) {
    var $username = $('#login-uname');
    var $usernameLbl = $('#login-uname-lbl');
    var $pwd= $('#login-pwd');
    var $pwdLbl = $('#login-pwd-lbl');
    var $h1 = $('#login-h1');
    var failed = false;
    $username.css('background-color', 'white');
    $pwd.css('background-color', 'white');
    $usernameLbl[0].firstChild.data = 'Username';
    $pwdLbl[0].firstChild.data = 'Password';
    $h1[0].firstChild.data = 'Login';
    if (username.length == 0) {
        $username.css('background-color', 'red');
        $usernameLbl[0].firstChild.data = 'Username - Cannot be empty';
        failed = true;
    }
    if (username.length > 20) {
        $username.css('background-color', 'red');
        $usernameLbl[0].firstChild.data = 'Username - Cannot be over 20 characters';
        failed = true;
    }
    if (password.length < 8) {
        $pwd.css('background-color', 'red');
        $pwdLbl[0].firstChild.data = 'Password - Must be at least 8 characters';
        failed = true;
    }
    if (failed) {
        return;
    }
    $.ajax({
        url: "/login/",
        data: {username: username, password: password}
    }).done(function(data) {
        data = JSON.parse(data);
        if (data.Error) {
            $username.css('background-color', 'red');
            $pwd.css('background-color', 'red');
            $h1[0].firstChild.data = 'Login - ' + data.Error;
        }
        if (data.Result) {
            location.reload();
        }
    });
};

oc.user.register = function(username, password, passwordDuplicate) {
    var $username = $('#reg-uname');
    var $usernameLbl = $('#reg-uname-lbl');
    var $pwd= $('#reg-pwd');
    var $pwd2 = $('#reg-pwd-dup');
    var $pwdLbl = $('#reg-pwd-lbl');
    var $pwd2Lbl = $('#reg-pwd-dup-lbl');
    var $h1 = $('#reg-h1');
    $username.css('background-color', 'white');
    $pwd.css('background-color', 'white');
    $pwd2.css('background-color', 'white');
    $usernameLbl[0].firstChild.data = 'Username';
    $pwdLbl[0].firstChild.data = 'Password';
    $pwd2Lbl[0].firstChild.data = 'Reenter Password';
    $h1.text('Register');
    var failed = false;
    if (username.length == 0) {
        $username.css('background-color', 'red');
        $usernameLbl[0].firstChild.data = 'Username - Cannot be empty';
        failed = true;
    }
    if (password.length < 8) {
        $pwd.css('background-color', 'red');
        $pwdLbl[0].firstChild.data = 'Password - Must be at least 8 characters';
        failed = true;
    }
    if (password !== passwordDuplicate) {
        $pwd2.css('background-color', 'red');
        $pwd2Lbl[0].firstChild.data = 'Reenter Password - Must be identical';
        failed = true;
    }
    if (failed) {
        return;
    }
    $.ajax({
      url: "/register/",
      data: {username: username, password: password}
    }).done(function(data) {
        data = JSON.parse(data);
        if (data.Error) {
            $username.css('background-color', 'red');
            $pwd.css('background-color', 'red');
            $pwd2.css('background-color', 'red');
            $h1[0].firstChild.data = 'Register - ' + data.Error;
        }
        if (data.Result) {
            location.reload();
        }
    });
};

oc.user.logout = function() {
    $.ajax({
        url: "/logout/"
    }).done(function(data) {
        data = JSON.parse(data);
        if (data.Error) {
            console.log(data.Error);
        }
        if (data.Result) {
            location.reload();
        }
    });
};
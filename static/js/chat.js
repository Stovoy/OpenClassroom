oc = oc || {};

oc.chat = oc.chat || {};

var identifier = Math.random();
var $dataStore = $('#data-store');
var id = $dataStore.data('page');

$(document).ready(function() {
    setInterval(oc.chat.refresh, 1000);
    oc.chat.refresh();

    $('#chat-input').on('keydown', function(e) {
        if (e.keyCode == 69) { // Enter
            oc.chat.sendMessage();
        }
    })
});

oc.chat.sendMessage = function() {
    
};

oc.chat.refresh = function() {
    $.ajax({
        url: "/chat/loadNew/",
        data: {
            page: id,
            lastMessage: -1,
            identifier: identifier}
    }).done(function(data) {
        data = JSON.parse(data);
        if (data.Error) {

        } else if (data.Users) {
            var users = data.Users;
            var $usersHeader = $('#users-header');
            var $chatUserArea = $('#chat-user-area');
            for (var i = 0; i < users.length; i++) {
                $usersHeader.siblings().remove();
                var color = "l2";
                if (i % 2 == 0) {
                    color = "l1";
                }
                var type = "guest-line";
                if (!users[i].IsGuest) {
                    type = "user-line";
                }
                $chatUserArea.append(
                        '<span class="' +
                        type + ' ' +
                        color + '">' +
                        users[i].Name +
                        '</span>');
            }
        }
    });
};
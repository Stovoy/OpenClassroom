oc = oc || {};

oc.chat = oc.chat || {};

var identifier = Math.random();
var $dataStore;
var id;
var $chatInput;
var $chatButton;
var $chatHistory;
var lastMessageID = 0;

$(document).ready(function() {
    // Load datastore
    $dataStore = $('#data-store');
    id = $dataStore.data('page');

    setInterval(oc.chat.refresh, 1000);
    oc.chat.refresh();

    $chatInput = $('#chat-input');
    $chatButton = $('#chat-button');

    $chatHistory = $('#chat-history');

    $chatInput.on('keydown', function(e) {
        if (e.keyCode == 13) { // Enter
            oc.chat.sendMessage();
        }
    });

    $chatButton.on('click', function() {
        oc.chat.sendMessage();
    })
});

oc.chat.sendMessage = function() {
    var message = $chatInput.val();
    $.ajax({
        method: "POST",
        url: "/chat/message/",
        data: {
            page: id,
            message: message}
    }).done(function(data) {
        data = JSON.parse(data);
        if (data.Error) {
            console.log(data);
        } else {
            $chatInput.val("");
            oc.chat.refresh();
        }
    });
};

oc.chat.refresh = function() {
    $.ajax({
        url: "/chat/loadNew/",
        data: {
            page: id,
            lastMessage: lastMessageID,
            identifier: identifier}
    }).done(function(data) {
        data = JSON.parse(data);
        if (data.Error) {
            console.log(data);
        }
        var i;
        if (data.Users) {
            var users = data.Users;
            var $usersHeader = $('#users-header');
            var $chatUserArea = $('#chat-user-area');
            for (i = 0; i < users.length; i++) {
                $usersHeader.siblings().remove();
                var color = "l2";
                if (i % 2 == 0) {
                    color = "l1";
                }
                var type = "guest-line";
                var name = users[i].Name;
                if (!users[i].IsGuest) {
                    type = "user-line";
                    name =
                        '<a href="/user/' + name +
                            '">' + name + '</a>'
                }
                $chatUserArea.append(
                        '<span class="' +
                        type + ' ' +
                        color + '">' +
                        name +
                        '</span>');
            }
        }
        if (data.NewMessages) {
            var messages = data.NewMessages;
            for (i = 0; i < messages.length; i++) {
                if (messages[i].ID > lastMessageID) {
                    lastMessageID = messages[i].ID;
                }
                $chatHistory.append(
                    '<div class="message">' +
                        '<a href="/user/' + messages[i].User + '">' +
                        messages[i].User + '</a>[' +
                        messages[i].Time + ']: ' +
                        messages[i].Message +
                    '</div>');
            }
        }
    });
};
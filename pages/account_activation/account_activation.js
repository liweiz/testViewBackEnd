// Validate inputted password
$('form').on('submit', function(){
  var passed = false 
  $('#my-btn').button('loading')
  if $('#password1').val().length > 7 && $('#password1').val().length < 21 {
    if($('#password1').val() == $('#password2').val()){
      passed = true
      var url = "http://localhost:3000/users/:user_id/password"
      var content = {};
      content["password"] = $('#password1').val();
      $.ajax(url, {
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(content)
      })
        .done(function() {
          $("#setting").addClass("hidden")
          $('#messageOK').removeClass("hidden")
          $('#messageOK').addClass("show")
        })
        .fail(jqXHR.fail(function(a, b, c) {
          $("#setting").addClass("hidden")
          $('#errorMessage').text(c)
          $('#messageFailed').removeClass("hidden")
          $('#messageFailed').addClass("show")
        }))
        .always(function() {
          $('#my-btn').button('reset')
        });
    }
  }
  if !passed {
    $('#my-btn').button('reset')
  }

   
});
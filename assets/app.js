// Form controllers
var registrationFormController = {
    submit: function(obj, event){
        event.preventDefault();
        NotificationService.clear();
        FormService.clear();
        var dataToServer = SerializeForm(obj);
        return Send("/api/v1/user/", "post", JSON.stringify(dataToServer), function (result) {
            if (result) {
                if (result.status === false){
                    FormService.showError(result.errors, obj)
                }else{
                    location.href = result.redirect
                }
            }
        });
    },
    googleRegistration: function () {

    }
};

var loginFormController = {
    submit: function(obj, event){
        event.preventDefault();
        NotificationService.clear();
        var dataToServer = SerializeForm(obj);
        return Send("/api/v1/session/", "post", JSON.stringify(dataToServer), function(result){
            if (result) {
                if (result.status === false){
                    NotificationService.showError(result.description)
                }else{
                    location.reload(true)
                }
            }
        });
    }
};

var profileFormController = {
    submit: function(obj, event){
        event.preventDefault();
        NotificationService.clear();
        FormService.clear();
        var dataToServer = SerializeForm(obj);
        return Send("/api/v1/user/", "put", JSON.stringify(dataToServer), function (result) {
            if (result) {
                if (result.status === false){
                    FormService.showError(result.errors, obj)
                }else{
                    location.href = result.redirect
                }
            }
        });
    }
};

var forgetPasswordController = {
    submit: function(obj, event){
        event.preventDefault();
        console.log("forget");
        NotificationService.clear();
        FormService.clear();
        var dataToServer = SerializeForm(obj);
        return Send("/api/v1/password/", "post", JSON.stringify(dataToServer), function (result) {
            if (result) {
                if (result.status === false){
                    FormService.showError(result.errors, obj)
                }else{
                    NotificationService.showInfo("Manual sent to your email.")
                }
            }
        });
    }
};

var newPasswordController = {
    submit: function(obj, event){
        event.preventDefault();
        NotificationService.clear();
        FormService.clear();
        var dataToServer = SerializeForm(obj);
        return Send("/api/v1/password/", "put", JSON.stringify(dataToServer), function (result) {
            if (result) {
                if (result.status === false){
                    FormService.showError(result.errors, obj)
                }else{
                    NotificationService.showInfo("New Password created. Use Sign In page.")
                }
            }
        });
    }
};

// Controller page
var PageController = {
    RegistrationTypeChange: function(obj, event){
        console.log("Changed", obj);
        event.preventDefault();
        var emailRegistrationFormContainer = document.getElementById("emailRegistrationForm");
        var typeRegistrationFormContainer = document.getElementById("typeRegistrationForm");
        typeRegistrationFormContainer.classList.add("d-none");
        emailRegistrationFormContainer.classList.remove("d-none");
    },
    EditProfile: function (obj, event) {
        event.preventDefault();
        console.log("print");
        var formUserProfile = document.getElementById("formUserProfile");
        var fields = document.getElementsByClassName("form-control");
        var labels = document.getElementsByClassName("form-value");
        var submitButton = document.getElementById("btnSaveProfile");
        // active form if it is hidden
        if (submitButton.classList.contains("d-none")){
            obj.textContent = "Cancel";
            for(var elm of fields){
                elm.classList.remove("d-none");
            }
            for(var elm of labels){
                elm.classList.add("d-none");
            }
            submitButton.classList.remove("d-none");
        }else{
            // chancel editing
            location.reload();
        }
    }
};

var GoogleMapController = {
    initAutocomplete: function () {
        var autocompleteField = document.getElementById('inputAddress');
        if (autocompleteField){
            autocomplete = new google.maps.places.Autocomplete(
                (autocompleteField), {types: ['geocode']});
        }
        autocomplete.addListener('place_changed', this.initMap);
    },
    initMap: function(){
        map = new google.maps.Map(document.getElementById('mapContainer'), {
            zoom: 8,
            center: new google.maps.LatLng(-34.397, 150.644)
        });
        var autocompleteField = document.getElementById('inputAddress');
        var geocoder = new google.maps.Geocoder();
        var address = autocompleteField.value;
        if (address !== ""){
            geocoder.geocode({'address': address}, function (results, status) {
                if (status === 'OK') {
                    map.setCenter(results[0].geometry.location);
                    var marker = new google.maps.Marker({
                        map: map,
                        position: results[0].geometry.location
                    });
                }
            });
        }
    }
};

/// Utilities
function SerializeForm(form) {
    var elements = form.elements;
    var dataContainer ={};
    for(var i = 0 ; i < elements.length ; i++){
        var item = elements.item(i);
        if (item.name !== "" && item.name !== undefined){
            dataContainer[item.name] = item.value;
        }
    }
    return dataContainer
}


function Send(url, method, body, handler) {
    fetch(url, {
        method: method,
        credentials: 'same-origin',
        headers: {
            'Accept': 'application/json, text/plain, */*',
            'Content-Type': 'application/json'
        },
        body: body
    }).then(function(response){return response.json();}).then(function(responseObject){
       handler(responseObject);
    }).catch(function(error){
        console.log(error);
    })
}

var NotificationService = {
    errorContainerId: "alertContainer",
    getNotificationTemplate: function(type, message){
        var notificationElement = document.createElement("div");
        notificationElement.setAttribute("class", "alert alert-"+type.toString());
        notificationElement.innerText = message;
        return notificationElement
    },
    showError: function(message){
        document.getElementById(this.errorContainerId).appendChild(this.getNotificationTemplate("danger", message));
    },
    showInfo: function(message){
        document.getElementById(this.errorContainerId).appendChild(this.getNotificationTemplate("info", message));
    },
    clear: function () {
        var container = document.getElementById(this.errorContainerId);
        while (container.firstChild) {
            container.removeChild(container.firstChild);
        }
    }
};

var FormService = {
    getErrorTemplate: function(message){
        var errorElement = document.createElement("div");
        errorElement.setAttribute("class", "invalid-feedback");
        errorElement.innerText = message;
        return errorElement
    },
    showError: function (listError, form) {
        var _this = this;
        console.log(listError);
        listError.forEach(function (value) {
            var element = form.querySelector("input[name="+value.name+"]");
            var errorElement = _this.getErrorTemplate(value.message);
            console.log(element);
            var parent = element.parentNode;
            var next = element.nextSibling;
            if (next){
                // console.log("Next", next, "parent", parent);
                parent.insertBefore(errorElement, next);
            }else{
                // console.log("Parent", next);
                parent.appendChild(errorElement);
            }
        });
    },
    clear: function (){
        var listError = document.getElementsByClassName("invalid-feedback");
        console.log(listError);
        if (listError) {
            for(var value of listError){
                value.remove();
            }
        }
    }
};
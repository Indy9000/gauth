import * as ko from "knockout";

class HelloViewModel{
    profileID:KnockoutObservable<string> = ko.observable("blank")
    profileName:KnockoutObservable<string> = ko.observable("blank")

    constructor(){
        gapi.signin2.render('my-signin2', {
            'scope': 'profile email',
            'width': 240,
            'height': 50,
            'longtitle': true,
            'theme': 'light',
            'onsuccess': param => this.onSignIn(param)
        });        
    }
    
    public onSignIn(googleUser:gapi.auth2.GoogleUser) {
        var profile = googleUser.getBasicProfile();
        this.profileID(profile.getId())
        this.profileName(profile.getName())

        console.log("ID: " + this.profileID); // Don't send this directly to your server!
        console.log('Full Name: ' + profile.getName());
        console.log('Given Name: ' + profile.getGivenName());
        console.log('Family Name: ' + profile.getFamilyName());
        console.log("Image URL: " + profile.getImageUrl());
        console.log("Email: " + profile.getEmail());

        // The ID token you need to pass to your backend:
        var id_token = googleUser.getAuthResponse().id_token;
        console.log("ID Token: " + id_token);
    };
}

ko.applyBindings(new HelloViewModel())

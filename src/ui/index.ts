import * as ko from "knockout";
import { AuthManager } from "./auth-manager";
import { UserProfile } from "./user-profile";

class MainViewModel{
    authMgr:KnockoutObservable<AuthManager>
    constructor(){
        let el = document.getElementById("Google-SignInButton")
        if (el){
            this.authMgr= ko.observable( new AuthManager(el) )
        }else{
            throw new Error('ERROR: HTML must contain an element with id `Google-SignInButton`');
        }
    }
}

ko.applyBindings(new MainViewModel())

/*
class HelloViewModel{
    userProfile:KnockoutObservable<UserProfile> = ko.observable(new UserProfile())

    constructor(){
        this.getProfile()
    }

    // Before authenticating, this tries to validate an existing session token
    // and retrieve the profile. If not successful, it presents a login button
    getProfile(){
        fetch('/api/user',{
            credentials:'include',
            method:'GET',
            headers:{
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            }
        })
        .then(response=>{
            if (response.status!=200){
                console.log("no session, authenticate first")
                this.renderGoogleSignInButton()
            }else{
                return response.json()
            }
        })
        .then(j=>{
            if (j){
                let p = this.makeProfie(j)
                this.userProfile(p)
            }
        })
    }

    renderGoogleSignInButton(){
        gapi.signin2.render('my-signin2', {
            'scope': 'profile email',
            'width': 100,
            'height': 30,
            'longtitle': false,
            'theme': 'light',
            'onsuccess': param => this.onSignIn(param)
        });        
    }

    public onSignIn(googleUser:gapi.auth2.GoogleUser) {
        var profile = googleUser.getBasicProfile();

        // The ID token you need to pass to your backend:
        var id_token = googleUser.getAuthResponse().id_token;
        console.log("ID Token: " + id_token);

        let p = new UserProfile()
        p.UserName = profile.getName()
        p.PhotoURL = profile.getImageUrl()
        p.UniqueUserID = profile.getId()

        this.postIDToken(p,id_token)
    }

    postIDToken(p:UserProfile,idToken:string){
        let payload = JSON.stringify(p)
        fetch('/api/user/auth',{
            method:'POST',
            body:payload,
            credentials:'include',
            headers:{
                'Content-Type':'application/json',
                'Authorization': "Bearer " + idToken
            }
        })
        .then(resp=>{
            if (resp.status==401){//unauthorised
                console.log("failed authentication. Something is a miss...")
                //TODO: show HTML for this failure
                return
            }
            return resp.json()
        })
        .then(j=>{
            console.log("success authentication")
            let p = this.makeProfie(j)
            this.userProfile(p)
        })
    }

    makeProfie(j:any):UserProfile{
        let p = new UserProfile()
        p.UserName = j["UserName"]
        p.PhotoURL =j["PhotoURL"]
        p.UniqueUserID =j["UniqueUserID"]
        console.log("Profile:",p)
        return p
    }
}

ko.applyBindings(new HelloViewModel())
*/
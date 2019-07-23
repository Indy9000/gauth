import * as ko from "knockout"
import { UserProfile } from "./user-profile";

export class AuthManager{
    gauth2!: gapi.auth2.GoogleAuth
    isGAuthInitialized:boolean = false
    isLoggedIn = ko.observable(false)
    isSignedUp = ko.observable(false)
    ctaCategory = ko.observable("")
    userProfile = ko.observable( new UserProfile() )
    cliConfig!: gapi.auth2.ClientConfig
    authToken = ""
    authErrorNotificication = ko.observable("")
    isOfflineMode = ko.observable(false)
    profileKey = "my-user-profile-v1"
    constructor(private elementToAttachClick:HTMLElement){

        // NOTE:    When the app is loaded, it first checks
        //          whether there's a stored profile in the localstorage
        // 
        //          if a profile is there, then that user had signed up before
        //              therefore go and log them in
        //          else
        //              show the sign up related content and consent
        //              show sign-in just in case user had signed up but cleared browser
        //                  thereby clearing the localstorage
        //              on sign up consent, signup endpoint is called and necessary 
        //                  table entries would be made
        // 
        //          In either case, google token is sent to the backend and backend
        //          verifies the token validity
        //          

        if (navigator.onLine){
            this.loadGAuthLib()
            this.isOfflineMode(false)
        }else{
            // set offline mode and try later
            this.isOfflineMode(true)
        }

        let signedUp =  this.signupCheck()
        this.isSignedUp(signedUp)
        if(signedUp){
            // TODO: attempt to authenticate current profiled user
            //       if the internet is not there, fall back to offline mode
            //          in this mode, app should work with the current profile stored
        }else{ //not signed up yet
            // TODO: Attempt to go into the landing and sign up/sign in stage
            //          this may fail due to internet not being there.
            //          In that case, we can't continue

            let isSignup = true // default to sign up
            this.AttachClickHandlerWithSelector(this.elementToAttachClick, isSignup)
        }
    
    }

    private loadGAuthLib(){
        this.cliConfig = {
            client_id: '56400899087-a1so0s71rmglqp0orpffs3bbr9atr66t.apps.googleusercontent.com',
            cookie_policy: 'single_host_origin',
            // Request scopes in addition to 'profile' and 'email'
            //scope: 'additional_scope'
        }
        gapi.load('auth2',()=>{
            this.gauth2 = gapi.auth2.init(this.cliConfig)
            this.isGAuthInitialized = true
        })
    }

    // CheckIfSignedIn checks if a user is signed in with the browser already
    // and if it is the user we have in the stored profile
    public CheckIfSignedIn(email:string){
        if (!this.isGAuthInitialized){
            return
        }

        // Is this user already signed in?
        if (this.gauth2.isSignedIn.get()) {
            var googleUser = this.gauth2.currentUser.get();
            
            console.log("googleUser.getBasicProfile():", googleUser.getBasicProfile())
            
            // Same user as in the credential object ?
            if (googleUser.getBasicProfile().getEmail() === email) {
                let uprof =  new UserProfile()
                uprof.UniqueUserID = googleUser.getBasicProfile().getId()
                uprof.UserName = googleUser.getBasicProfile().getName()
                uprof.PhotoURL = googleUser.getBasicProfile().getImageUrl()
        
                this.setLoggedInUser(uprof)
            }
        }

        // TODO: Here a flag should be set to indicate that the user is not
        //       signed in. This should be bound to a pop-up of some sort which
        //       gets user input and signs in. Once the sign-in had completed
        //       background syncing can continue.
        this.isLoggedIn(false)
    }

    private setLoggedInUser(uprof:UserProfile){
        this.userProfile(uprof)
        this.isLoggedIn(true)
        //TODO: set the profile to localstorage
        let pj = JSON.stringify(uprof)
        localStorage.setItem(this.profileKey, pj)
        this.isSignedUp(true)
        console.log("setLoggedInUser",uprof)
    }

    public AttachClickHandlerWithoutSelector(element:HTMLElement, isSignup:boolean){
        let options: gapi.auth2.SigninOptions = {
            fetch_basic_profile:true
        }
        this.attachClickHandler(element, options, isSignup)
        console.log("AttachClickHandlerWithoutSelector")
    }

    public AttachClickHandlerWithSelector(element:HTMLElement, isSignup:boolean){
        let options: gapi.auth2.SigninOptions = {
            prompt:'select_account',
            fetch_basic_profile:true        
        }
        this.attachClickHandler(element,options,isSignup)
    }

    private postUserProfile(p:UserProfile, endpoint:string){
        let payload = JSON.stringify(p)

        // console.log("postIDToken(),payload:",payload," authToken:",this.authToken)
        this.authErrorNotificication("")

        fetch(endpoint,{
            method:'POST',
            body:payload,
            credentials:'include',
            headers:{
                'Content-Type':'application/json',
                'Authorization': "Bearer " + this.authToken
            }
        })
        .then(resp=>{
            if (resp.status==401){//unauthorised
                let msg = "Authentication failed. Something is a miss..."
                this.authErrorNotificication(msg)
                console.log(msg,resp)
                return
            }
            return resp.json()
        })
        .then(j=>{
            if(j){
                console.log("postIDToken success authentication")
                // callback(j)
                this.setLoggedInUser(p)

            }else{
                // reloadCallback()
            }
        })
        .catch(err=>{
            console.log("Error while authenticating..", err)
        })
    }

    private attachClickHandler(element:HTMLElement,options: gapi.auth2.SigninOptions, isSignup:boolean){
        console.log(element.id);

        let attacher = ()=>{
            this.gauth2.attachClickHandler(element, options,
            (googleUser)=>{ // success
                // The ID token is passed to the backend for verification
                this.authToken = googleUser.getAuthResponse().id_token;
                // post the
                let p = new UserProfile()
                p.UserName = googleUser.getBasicProfile().getName()
                p.PhotoURL = googleUser.getBasicProfile().getImageUrl()
                p.UniqueUserID = googleUser.getBasicProfile().getId()

                let endpoint = '/api/v1/user/auth'
                if (this.ctaCategory() == "Sign up"){
                    // call backend for sign up, validate auth token
                    // validate UID, if user is already available just log them in
                    // if token can't be validated then return error
                    endpoint = '/api/v1/user/signup'
                }
                console.log("Endpoint selected:",endpoint)
                this.postUserProfile(p, endpoint)
            }, 
            function(error) {//failure to Sign in with 
                alert(JSON.stringify(error, undefined, 2));
            }
            )
        }

        if (this.isGAuthInitialized){
            attacher()
        }else{
            gapi.load('auth2',()=>{
                this.gauth2 = gapi.auth2.init(this.cliConfig)
                this.isGAuthInitialized = true
                attacher()
            })
        }
    }

    private signupCheck():boolean{
        // Signup check
        // If a user hasn't signed up, they should see the landing
        // page. Landing page should have signup form which
        // creates a user on the backend with associated google
        // unique userid
        // 
        // check if localStorage has profile details
        let profile_value = localStorage.getItem(this.profileKey)
        if (profile_value){
            // TODO: deserialize to proflle
            let obj =  JSON.parse(profile_value)
            let p =   new UserProfile()
            p.UniqueUserID = obj["UniqueUserID"]
            p.PhotoURL = obj["PhotoURL"]
            p.UserName = obj["UserName"]
            this.userProfile(p)
            console.log('User signed up already')
            return true
        }else{// redirect to Landing
            console.log('User is not signed up yet')
            // window.location.href = 'landing.html'
            return false
        }
    }

    public SetCTACategory():boolean{
        console.log(this.ctaCategory())
        return true
    }
    // // NOTE: signIn pop-ups a dialog to choose the identity. Must be called from a click
    // // handler
    // private signIn(){
    //     console.log("user not signed in")
    //     let options = new gapi.auth2.SigninOptionsBuilder()
    //     options.setPrompt('select_account')
    //     // Otherwise, run a new authentication flow.
    //     this.gauth2.signIn(options)
    //     .then((googleUser)=>{
    //         console.log("signed in..")
    //         console.log("googleUser.getBasicProfile():", googleUser.getBasicProfile())
    //         this.isLoggedIn(true)
    //     })
    //     .catch((err)=>{
    //         this.isLoggedIn(false)
    //         console.log("error:",err)
    //         if( err.error =="popup_blocked_by_browser"){
    //             console.log("Log in selector was blocked by browser")
    //         }
    //     })
    // }

}
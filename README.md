# OutaTime
OutaTime Events is an event-hosting application where users can create accounts to host events. Each user will have their own page with organization/host description, links to social media, images, descriptions and more. Their page will be customizable, and each organization will have a URL they can give out or place on their website/social media.

Events can be made private to be shown only to desired people or public, where people can search for and view the event. Every event will have a time, description, and a location, helping attendees find their way.

Those visiting the app can subscribe to users and events to get notifications about new events, upcoming events, changes etc., no need for an account, only an email. With a unique code, subscribers can view their subscriptions and cancel them as desired.

## Run

In the main root folder, run
```bash
go run main.go
```

The app is running on localhost:3000

## Random String Utility

The random string utility will return a string made up of random letters and numbers.

It will create a string of x groups, made up of n characters, separated by a dash. (5, 5) Proves to be unique enough. Tests show no repeats in 100,000 generations. 

This utility will be used to create unique and random authentication codes for users, which will be saved to the database and as a cookie in the browser. This way, a user can remain logged in without having to worry about other people guessing their value, as it is a long, complex, and random string.

This same utility will also be used to create nonce values in forms to prevent cross-site request forgeries, or “CSRF attacks”.

This will also be used to create codes for private events to be used in their routing.

## Pages Design

The Go templates will be made using HTML.

All pages will have the same basic design and layout. This will be created in the layout.html page. It will contain a reference to a `{{template “body” .}}`, which will contain the specific structure of any given page, as well as pass on any data it needs to be filled out.

All pages will be made with `{{define “body”}}` and `{{end}}` tags. This way, they can be integrated into the layout.html template.

When the pages are rendered, they will require the layout.html as well as the page containing the “body” template. The data passed into the template will be dependant on which page is being rendered.

## New Accounts

When a user creates an account, they will enter hostname/organization, which will be displayed on their homepage, as well as on their events. The database will store the username given, as well as a “URL-friendly” name with spaces and special characters removed. The given username will be shown to users, and the URL-friendly name will be used in routing.

Users may link their website and social media links.

After signing up, a verification link will be sent to the email provided when signing up. The account made on the database will be invalid and inactive until verification has been done, after which, it will be valid and active.

Users will be able to customize their own page. Their landing pages will be a long form with fields for “about us”, change avatar, change page banner, and they will be able to see all events they have made. They can edit, delete, and add new events. 

The user’s page will be unpublished until the user pushes the “publish” button at the bottom. At any time, the user can unpublish the page. An unpublished page will be hidden from the public, along with all events. When making the page, the user will have the ability to view page as a viewer, by a button.

When signing up, the password given will be hashed for security.

## User Authentication

When a user signs up, the password will be hashed by JavaScript and then sent to the sever to be checked against the database which will also contain a hashed version of the password. If there is a match, the user will be logged in.

An authentication code will be generated by the Random String utility and saved in the database as well as in a cookie in the browser. When logging out, these will both be removed.

When the user accesses the app, if an authentication key is in the browser, it will be checked against the database, and the user data will fill out the templates. This check will be done in middleware.

## Event Creation

When making an event, users will be able to set a date and time, which will be displayed. The day before the event is to take place, notifications will be sent out to subscribers.

Events can be listed as private or public. Public events will be visible on the user’s page and routes and may appear on the OutaTime Events home page. Viewers of the app will be able to search for it as well. The route for these pages will be made using the event ID.

Private events will not be on any page. The link to it will be seen only by the user and to whom they choose to give the link. Private events will have a random string assigned to them (event code), which will be used in the routing. This way, the user can still give out a link, but the general public will have trouble guessing the URL, as it is a long sequence of random characters.

Users will also be able to set a location using a Google Maps API and set an image for it.

Users will enter two descriptions: a short one to be used on the “card” and a longer one on the event page.
Events will be seen as a card, showing the image, host, time and short description. Clicking on the card will bring the viewer to the event page, where they will see the image, the host (a link to the host’s home page), the time, the full description and the Google Maps location.

When creating an event, the user will have the ability to view the card and view the event page as a viewer, by a button.

## Subscriptions

Viewers will be able to subscribe to either an event or to an organization itself.

Subscriptions to an organization will mean that subscribers will get notifications about any and all of the events that the organization has.

Subscriptions to an event will only result in notifications about that specific event. After which, the subscription, along with the event, will be removed from the database.

Viewers do not need to make an account to subscribe. They need only to enter their email. A verification link will be sent to verify the email.

Those subscribed to an organization will not be able to resubscribe, nor subscribe to any of that organization’s events. Trying this will return, “that email is already subscribed”.

Subscribers will be given a subscriber code, which they can use to check their subscriptions and unsubscribe as desired on the “my subscriptions” page.

There will be a separate table for subscribers and their codes, as well as a valid field, and a field for what they are subscribed to.

## Administration

The “admin” account will be added manually to the database. There will be only **one** admin account. All other accounts will be normal by default.

The admin account will be able to view all users and their events, images, data, etc. Admin will have the ability to delete any event, as well as suspend or remove any user. If admin deletes an event, or user, a notice will be sent out to subscribers, along with an explanation.

## Daily Function

A function will run once a day to send emails to subscribers about events happening the next day.
It will also check all emails that have not been verified within the last day and remove them.


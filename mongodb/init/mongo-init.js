init = false;
print("Init script ...")

try {
  if (!db.isMaster().ismaster) {
    print("Error: primary not ready, initialize ...")
    rs.initiate(
      {
        _id:'my-replica-set',
        members: [
          { _id:0,
            host: "mongo:27017"
          }
        ]
      }
    )
    quit(1);
  } else {
    if (!init) {
      admin = db.getSiblingDB("admin");
      admin.createUser(
        {
          user: "root",
          pwd: "password",
          roles: ["readWriteAnyDatabase"]
        }
      );

      db = db.getSiblingDB("purchase");
      db.createUser({
        user: "editor",
        pwd: "Tmhr1582",
        roles: [
          { role : "dbOwner", db : "purchase" },
          { role : "dbAdmin", db : "purchase" },
          { role : "readWrite", db : "purchase" },
        ]
      });

      db.createCollection("subscriptionOrders");
      db.subscriptionOrders.insertMany([
        {
          _id:"1456591d:01c637d0-c225-4b39-83a5-06f39fbf015e",
          userId:"9gab33rqtgptMd",
          planCode:"subscription_premium",
          status: 2,
          automaticRenewal: true,
          periodEndAt: 1719682304,
          trialEndAt: 0,
          gracePeriod: 0,
          purchaseType: 1,
          productId:"subscription_premium",
          activatedAt: 1698600710,
          updatedAt: 1716975125,
          createdAt: 1698600710,
          statusDetail: 0
        }
      ]);

      init = true;
    }
  }
} catch(e) {
  rs.status().ok
}

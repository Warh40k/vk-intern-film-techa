BEGIN;

ALTER table films ALTER column released DROP NOT NULL;
ALTER table films ALTER column description DROP NOT NULL;
ALTER table films ALTER column rating DROP NOT NULL;

ALTER table actors ALTER column birthday DROP NOT NULL;
ALTER table actors ALTER column gender DROP NOT NULL;

END;
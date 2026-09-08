# Legal & Compliance Specification (India IT Rules 2021) - TN NOW

As an intermediary, TN NOW strictly complies with India's **Information Technology (Intermediary Guidelines and Digital Media Ethics Code) Rules, 2021**.

---

## 1. Designated Grievance Redressal Officer

To comply with Rule 3(2), TN NOW publishes the details of the Grievance Officer in-app (Settings / Support) and on a public web page.

- **Name**: Grievance Redressal Officer, TN NOW
- **Email**: compliance@tnnow.in
- **Mailing Address**: TN NOW Compliance Desk, Madurai, Tamil Nadu, India.

---

## 2. Grievance Tracking & SLA Workflow

All grievance complaints are recorded in the `grievances` database table.

```text
User Submits Grievance -> Database Log -> Auto-Ack Email (SLA: 24 Hours) -> Moderator Review -> Action & Reply (SLA: 15 Days)
```

1. **Acknowledgment**: Sent automatically within **24 hours** with a unique tracking ID.
2. **Resolution SLA**: Complete investigation and dispatch response within **15 days** of receipt.
3. **Emergency Takedowns (Sexual Content)**: If a complaint is received regarding content depicting nudity, sexual acts, or impersonation, content must be removed/disabled within **24 hours** (Rule 3(2)(b)).

---

## 3. Lawful Takedown Orders Compliance

Upon receipt of a court order or direction from a government agency authorized by law:
- TN NOW will remove or disable access to the specified content within **36 hours** of receiving the order (Rule 3(1)(d)).
- Action is logged in the `audit_logs` table for compliance tracking.

---

## 4. User Rights Controls

1. **User Blocking**: Every user can block another user. The blocked user's content, comments, and profile will be completely hidden from the blocker's view.
2. **Account Deletion**: Users can request account and data deletion in Settings. Deletion anonymizes or deletes user rows, profile media, and session logs from primary databases.

---

## 5. Significant Social Media Intermediary (SSMI) Tracker

If the registered user base of TN NOW crosses **5 million users**:
- Additional compliance steps trigger, including publishing monthly compliance reports, appointing a Chief Compliance Officer, a Nodal Contact Person, and a Resident Grievance Officer.
- An alert is configured in the Admin Dashboard when registered users cross **4.5 million** to prepare compliance onboarding.

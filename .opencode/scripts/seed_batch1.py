import urllib.request, json, sys

KEY = 'ntn_z13943653061c5sQPOZJLETAERl4sUiX4pzTUvF5eGebtA'
DB  = '3250531c-537a-8109-aa3d-cd3961fe0257'
H   = {
    'Authorization': 'Bearer ' + KEY,
    'Notion-Version': '2022-06-28',
    'Content-Type': 'application/json',
}

def ins(tid, story, pri, pts, tags):
    body = {
        'parent': {'database_id': DB},
        'properties': {
            'Task ID':  {'title':       [{'text': {'content': tid}}]},
            'Story':    {'rich_text':   [{'text': {'content': story}}]},
            'Priority': {'select':      {'name': pri}},
            'Points':   {'number':      pts},
            'Tag':      {'multi_select':[{'name': t} for t in tags]},
            'Status':   {'select':      {'name': 'BACKLOG'}},
        },
    }
    req = urllib.request.Request(
        'https://api.notion.com/v1/pages',
        data=json.dumps(body).encode(),
        headers=H,
        method='POST',
    )
    try:
        d = json.loads(urllib.request.urlopen(req, timeout=15).read())
        print('OK  ' + tid)
    except urllib.error.HTTPError as e:
        err = json.loads(e.read())
        print('ERR ' + tid + ' ' + err.get('message', str(e)))
    except Exception as e:
        print('ERR ' + tid + ' ' + str(e))

STORIES = [
    ('US-M01-01','As an Admin I want to log in with username/password so that I can access the system.','Must Have',3,['frontend','backend']),
    ('US-M01-02','As an Admin I want to be forced to change my password on first login so that default credentials are replaced.','Must Have',3,['frontend','backend']),
    ('US-M01-03','As an Admin I want my session to expire after 30 min of inactivity so that the system stays secure.','Must Have',2,['backend']),
    ('US-M01-04','As an Admin I want to view all system users so that I can manage accounts.','Must Have',2,['frontend','backend']),
    ('US-M01-05','As an Admin I want to create new HR/Supervisor accounts so that staff can access the system.','Must Have',3,['frontend','backend']),
    ('US-M01-06','As an Admin I want to edit and delete user accounts so that I can manage access rights.','Must Have',3,['frontend','backend']),
    ('US-M01-07','As an Admin I want to reset a user password directly so that locked-out users regain access immediately.','Must Have',2,['frontend','backend']),
    ('US-M02-01','As an HR/Supervisor I want to log in with username or email + password so that I can access the ATS.','Must Have',2,['frontend','backend']),
    ('US-M02-02','As an HR user I want to reset my password via OTP sent to my email so that I can regain access without Admin help.','Must Have',5,['frontend','backend']),
    ('US-M03-01','As an Applicant I want to open a job posting via a direct URL without registering so that I can apply easily.','Must Have',3,['frontend','backend']),
    ('US-M03-02','As an Applicant I want to fill and submit an application form so that my candidacy is registered.','Must Have',5,['frontend','backend']),
    ('US-M03-03','As the system I want to block form submission from blacklisted email addresses so that banned candidates cannot re-apply.','Must Have',3,['backend']),
    ('US-M03-04','As the system I want to auto-trigger OCR on uploaded resumes so that HR gets extracted text for AI scoring.','Must Have',3,['backend','ai-service']),
    ('US-M03-05','As an Applicant I want my online test to auto-save and remain re-openable until I submit or it expires so that I do not lose progress.','Must Have',5,['frontend','backend']),
    ('US-M04-01','As HR I want to create a job posting with hard skill weights so that AI similarity scoring is calibrated.','Must Have',5,['frontend','backend']),
    ('US-M04-02','As HR I want jobs to auto-publish and auto-close on scheduled dates so that I do not manage it manually.','Must Have',5,['backend']),
    ('US-M04-03','As HR I want to edit, filter, and manage job postings so that the job list stays accurate.','Must Have',3,['frontend','backend']),
    ('US-M04-04','As HR I want a pipeline dashboard with a 5-stage milestone view so that I can track all applicants at a glance.','Must Have',5,['frontend','backend']),
    ('US-M04-05','As HR I want dashboard stats showing total, remaining, dropped, and hired counts so that I can monitor recruitment health.','Must Have',3,['frontend','backend']),
    ('US-M04-06','As HR I want to trigger email notifications from the pipeline dashboard so that applicants are informed at each stage.','Must Have',5,['frontend','backend']),
    ('US-M05-01','As HR I want to build application forms with standard field types so that applicants can submit structured data.','Must Have',5,['frontend','backend']),
    ('US-M05-02','As HR I want to add custom question types to forms so that I can collect role-specific information.','Must Have',5,['frontend','backend']),
    ('US-M05-03','As HR I want to publish a form as JSON to the database so that the Applicant Portal can render it dynamically.','Must Have',3,['backend']),
    ('US-M05-04','As HR I want to edit and delete forms with a reference check so that active jobs are not accidentally broken.','Must Have',3,['frontend','backend']),
]

print('Inserting %d Batch 1 stories...' % len(STORIES))
for s in STORIES:
    ins(*s)
print('Done.')

# taskDumb-V2

Day 12 Kurang Function Edit Button


INSERT INTO public.tb_project(
	project_name, description, image)
	VALUES ('Title ke-1', 'Isi dari title yang ke-1', 'project1.png');

INSERT INTO public.tb_project(
	project_name, description, image, start_date, end_date, technology)
	VALUES ('Title ke-4', 'Isi dari title yang ke-4', 'project4.png', '2023-11-01','2023-12-10',ARRAY['go','pyt','fig']);
	
-- SELECT ALL DATA FROM tb_project
SELECT id, project_name, description, image, post_date, author_id
	FROM public.tb_project;
	
SELECT * FROM tb_project;

UPDATE public.tb_project
    SET id=?, project_name=?, description=?, image=?, post_date=?
    SET description='TITLE BERUBAH'
    WHERE id=2

DELETE FROM public.tb_project
    WHERE id=?